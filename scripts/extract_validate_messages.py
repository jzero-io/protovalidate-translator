import re
import json
import os

def main():
    base = os.path.dirname(os.path.abspath(__file__))
    proto_path = os.path.join(base, "../third_party/buf/validate/validate.proto")
    with open(proto_path, "r", encoding="utf-8") as f:
        content = f.read()

    # 匹配 (predefined).cel = { ... } 块，考虑嵌套
    blocks = []
    i = 0
    while True:
        start = content.find("(predefined).cel = {", i)
        if start == -1:
            break
        depth = 1
        pos = content.index("{", start) + 1
        while depth > 0 and pos < len(content):
            if content[pos] == "{":
                depth += 1
            elif content[pos] == "}":
                depth -= 1
            pos += 1
        block = content[start:pos]
        blocks.append(block)
        i = pos

    result = {}
    for block in blocks:
        # 提取 id
        id_m = re.search(r'id:\s*"([^"]+)"', block)
        if not id_m:
            continue
        rid = id_m.group(1)

        # 优先使用 message
        msg_m = re.search(r'message:\s*"([^"]*)"', block)
        if msg_m:
            result[rid] = msg_m.group(1)
            continue

        # 从 expression 提取: 找 '...'.format 或 ? '...' : '' 中的字符串
        expr = block
        # 合并多行字符串（expression: "a" "b" 形式），并压缩空白便于正则匹配
        expr = re.sub(r'"\s*\n\s*"', " ", expr)
        expr = re.sub(r"\s+", " ", expr)

        msg = None
        # 模式1: 从 '.format([' 向前找消息：结束引号在 .format 前，起始引号在 ? ' 后
        fmt_pos = expr.find("'.format([")
        if fmt_pos > 0:
            end_quote = fmt_pos  # 结束引号就是 '.format 的 '
            # 在 [0, end_quote) 内找 ? '，其后的内容到 end_quote 即为消息
            start_m = re.search(r"\?\s*'", expr[:end_quote])
            if start_m:
                start_pos = start_m.end()
                msg = expr[start_pos:end_quote]
                if "value" in msg or "must" in msg or "map" in msg or "repeated" in msg:
                    msg = msg.replace("\\'", "'")
                else:
                    msg = None
        if not msg:
            # 模式2: ? '...' : ''（无 format，如 value must be finite）
            m = re.search(r"\?\s*'([^']*(?:\\'[^']*)*)'\s*:\s*''", expr)
            if m:
                msg = m.group(1).replace("\\'", "'")
        if not msg:
            # 模式3: 单行 expression 中的 ? '...' : ''
            m = re.search(r"\?\s*'([^']*)'\s*:\s*''", expr)
            if m:
                msg = m.group(1)

        if msg:
            # %s / %d / %x 替换为 {{.Value}}
            msg = re.sub(r"%[sdx]", "{{.Value}}", msg)
            result[rid] = msg.replace("`", "")
        else:
            print(f"no message found for {rid}")
        # 无 message 的（如 example）不写入，保持输出只含有效文案

    print(f"found {len(result)} messages")
    # 保存到文件（go-i18n JSON 格式）
    items = [{"id": k, "translation": v} for k, v in sorted(result.items())]
    out_path = os.path.join(base, "../translator/locales/en.json")
    with open(out_path, "w", encoding="utf-8") as f:
        json.dump(items, f, indent=2, ensure_ascii=False)

if __name__ == "__main__":
    main()
