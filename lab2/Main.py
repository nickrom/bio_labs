import NW
import NWAffine
from sys import argv
from Util import input_seq

if __name__ == '__main__':
    sequences = input_seq(argv[1])

    while True:
        print("""Введите номер алгоритма:
		1. Needleman-Wunsch
		2. Needleman-Wunsch Affine
		0. Exit""")

        alg_number = input()

        if alg_number == '1':
            NW.calculate(sequences[0], sequences[1])
        elif alg_number == '2':
            NWAffine.calculate(sequences[0], sequences[1])
        elif alg_number == '0':
            break
        else:
            raise Exception("Ожидается 1, 2 или 3")
        print("-" * 100)
