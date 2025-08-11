# Microlist---proggraming-language
an 8-bit abstract proggraming language

## only linux (or maybe not i dont know)

# how to download
first download the zip file and extract it
(if you havent already) install go 1.19 from https://go.dev/doc/install
move the folder to your home directory
then go to the extracted folder in the terminal and open it run the command 
```
go build -o microlist main.go
```
if it says there is no main.go type "ls" and then move to that folder using cd
them move microlist to your home directory
if you have done that move it to your bin directory
```
sudo mv ~/microlist /bin/

```

YOU INSTALLED MICROLIST!!!!!


# tutoriol for microlist

there is only one command in microlist it is....
```
ml
```

there 2 types of ml, ml that give output and ml's that dont give output
no output:
```
ml
```
with output:
```
ml.
```

lets make a simple program!
make a file with the .ml extention and open it...
in every ml only can be 8 bits stored but that is no problem!
lets make hello, world
```
1 ml.
"hello"
"world"
```
jipee you made your first program!!!
run it with:
```
microlist <filename>.ml
```
lets see if it workes!

now its getting more complicated... but hold on
see in the hello world program we made you see we did this
```
1 ml.
```
why the 1?
well the 1 is the name of the ml this is why:
we can use "name of ml[]" to get ml's in other ml's
IMPORTANT an ml without a dot is like a variable it only stores info
```
1 ml
"hello"

2 ml.
1[]
"world"

```
what you think this will do.... yes!! its prints "hello, world"


...
wait thats the whole language??

jes it is!!!!

i dont know where it could be used for... maybe it can run doom or something but you can figure it out!!

BYE!!!!!
