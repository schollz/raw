import json
import math
import random

import toml


def split(word):
    return [char for char in word]


class S:
    def __init__(self):
        self.max_bars = 80
        self.sections = []
        self.name_to_id = {}

    def add(self, name, abc, sync=None):
        self.name_to_id[name] = len(self.sections)
        abc_arr = split(abc)
        abc_len = len(abc_arr)
        section = []
        for i, part in enumerate(split(abc)):
            section.append(
                {"part": part, "pos": math.floor(self.max_bars * i / abc_len)}
            )
        self.sections.append({"parts": abc, "starts": section})
        if sync != None:
            self.sync(self.name_to_id[name], self.name_to_id[sync])

    def sync(self, x, y):
        for _, v in enumerate(self.sections[y]["starts"]):
            if v["pos"] == 0:
                continue
            if v["part"] in self.sections[x]["parts"]:
                self.sections[x]["starts"].append({"part": v["part"], "pos": v["pos"]})

    def eval(self):
        sections = []
        for _, section in enumerate(self.sections):
            evaled = ["0"] * self.max_bars
            for _, v in enumerate(section["starts"]):
                startval = ""
                for i in range(self.max_bars):
                    if i >= v["pos"]:
                        if startval == "":
                            startval = evaled[i]
                        elif startval != evaled[i]:
                            continue
                        evaled[i] = v["part"]
            sections.append(evaled)
        return sections

    def lengths(self):
        sections = []
        for section in self.eval():
            lengths = []
            length = 1
            v0 = section[0]
            for i, v in enumerate(section):
                if i == 0:
                    continue
                if v != v0:
                    lengths.append({"part": v0, "length": length})
                    length = 0
                    v0 = v
                length += 1
            lengths.append({"part": v0, "length": length})
            sections.append(lengths)
        return sections


s = S()
s.add("c", "ABABC")
s.add("d", "DEFDEFGAH", "c")
s.add("e", "BBACCAAACCCAC", "c")
for _, section in enumerate(s.eval()):
    print("".join(section))

print(json.dumps(s.lengths(), indent=2))

import toml

toml_string = """
[c]
parts="ABABC"
[c.A]
something=3
    
[d]
parts="DEFDEFGAH"
sync="c"
[d.D]
something=4

"""

parsed_toml = toml.loads(toml_string)
print(json.dumps(parsed_toml, indent=2))
