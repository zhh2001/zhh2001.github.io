SET age 20
# OK
OBJECT ENCODING age
# "int"
SET name zhang
# OK
OBJECT ENCODING name
# "embstr"
SET text aaaaabbbbbcccccdddddeeeeefffffggggghhhhhiiiii
# OK
OBJECT ENCODING text
# "raw"