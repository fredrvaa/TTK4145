# Python 3.3.3 and 2.7.6
# python helloworld_python.py

from threading import Thread

i = 0

def thread1():
    global i
    for j in range(1000000):
	i += 1

def thread2():
    global i
    for j in range(1000000):
	i -= 1


def main():
    someThread1 = Thread(target = thread1, args = (),)
    someThread2 = Thread(target = thread2, args = (),)
    someThread1.start()
    someThread2.start()
    
    someThread1.join()
    someThread2.join()
    print(i)


main()
