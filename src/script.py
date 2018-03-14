#!/usr/bin/env python

newlines = ""
with open('./tangov2.ebnf') as f:
    processing = False
    current = ''
    for line in f:
        newlines += line[:-1]
        if line.strip() == "<< import \"tango/src/ast\" >>":
            processing = True
            newlines += "\n"
            continue
        if not processing:
            newlines += "\n"
            continue
        if line.strip() == "":
            newlines += "\n"
            continue
        words = line.split()
        if words[0] == "//" or words[0] == "/*":
            newlines += "\n"
            continue
        if current == '':
            if words[1] == ':':
                current = words[0].strip()
            else:
                print(newlines)
                raise ValueError("Current Not Found")
            # Do something
            if words[2] == "empty":
                newlines += "\t<< ast.AddNode(\"" + current + "\") >>"
                # print(current, 0)
            else:
                newlines += "\t<< ast.AddNode(\"" + current + "\"" + ''.join([', $' + str(i) for i in range(len(words[2:]))]) + ") >>"
                # print(current, len(words[2:]))
        else:
            if words[0] == ";":
                current = ""
                newlines += "\n"
                continue
            if words[0] != "|":
                print(newlines)
                raise ValueError("Pipe Not Found")
            newlines += "\t<< ast.AddNode(\"" + current + "\"" + ''.join([', $' + str(i) for i in range(len(words[1:]))]) + ") >>"
            # print(current, len(words[1:]))
            # Do something
        newlines += "\n"

with open('./tango-main.ebnf', 'w') as f:
    f.write(newlines)
