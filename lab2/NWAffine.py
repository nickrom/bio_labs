import numpy

import BlosumReader
import Util


# задаем параметры
class Options:
    def __init__(self, seq_1, seq_2, subs_matrix = "blosum62.txt"):
        self.GAP_START = -10
        self.GAP_EXTEND = -1
        self.LEN_SEQ_1 = len(seq_1) + 1
        self.LEN_SEQ_2 = len(seq_2) + 1
        self.MIN = -float("inf")
        self.BLOSUM_MATRIX, self.LETTER_DICT = BlosumReader.load_matrix("./subs_matrices/" + subs_matrix)
        self.NAME_MATRIX = subs_matrix

    # сравниваем две буквы
    def compare(self, letter_1, letter_2):
        index_i = self.LETTER_DICT[letter_1]
        index_j = self.LETTER_DICT[letter_2]
        return self.BLOSUM_MATRIX[index_i][index_j]

    def info(self):
        print("Needleman-Wunsch algorithm with affine gaps.\n\
               \t Substitution matrix:" + self.NAME_MATRIX + "\n\
               \t Gap_start:" + str(self.GAP_START) + "\n\
               \t Gap_extend:" + str(self.GAP_EXTEND))

# определим метод для инициализации матриц
def init_matrix(constant):
    matrix_M = numpy.zeros(shape=(constant.LEN_SEQ_1, constant.LEN_SEQ_2))
    matrix_M[0, :] = constant.MIN
    matrix_M[:, 0] = constant.MIN
    matrix_M[0, 0] = 0

    matrix_I = numpy.zeros(shape=(constant.LEN_SEQ_1, constant.LEN_SEQ_2))
    matrix_I[:, 0] = [constant.GAP_START + constant.GAP_EXTEND * (j - 1) for j in range(constant.LEN_SEQ_1)]
    matrix_I[0, :] = constant.MIN

    matrix_D = numpy.zeros(shape=(constant.LEN_SEQ_1, constant.LEN_SEQ_2))
    matrix_D[0, :] = [constant.GAP_START + constant.GAP_EXTEND * (j - 1) for j in range(constant.LEN_SEQ_2)]
    matrix_D[:, 0] = constant.MIN

    return matrix_M, matrix_I, matrix_D


def calculate(seq_1, seq_2):
    options = Options(seq_1, seq_2)
    options.info()
    print("\nCalculating...\n")

    matrix_M, matrix_I, matrix_D = init_matrix(options)

    # заполним матрицы значениями
    for i in range(1, options.LEN_SEQ_1):
        for j in range(1, options.LEN_SEQ_2):
            matrix_D[i][j] = max(options.GAP_START + matrix_M[i - 1][j],
                         options.GAP_EXTEND + matrix_D[i - 1][j],
                         options.GAP_START + matrix_I[i - 1][j])

            matrix_I[i][j] = max(options.GAP_START + matrix_M[i][j - 1],
                         options.GAP_START + matrix_D[i][j - 1],
                         options.GAP_EXTEND + matrix_I[i][j - 1])

            matrix_M[i][j] = max(matrix_M[i - 1, j - 1],
                         matrix_D[i - 1, j - 1],
                         matrix_I[i - 1, j - 1]) + options.compare(seq_1[i - 1], seq_2[j - 1])

    # начинаем идти с нижней правой ячейки
    align1, align2 = '', ''
    i, j = len(seq_1), len(seq_2)
    score = matrix_M[i][j]
    while i > 0 and j > 0:
        max_val = max(matrix_M[i, j], matrix_I[i, j], matrix_D[i, j])
        if max_val == matrix_I[i][j]:
            align1 += '-'
            align2 += seq_2[j - 1]
            j -= 1
        elif max_val == matrix_D[i][j]:
            align1 += seq_1[i - 1]
            align2 += '-'
            i -= 1
        elif max_val == matrix_M[i][j]:
            align1 += seq_1[i - 1]
            align2 += seq_2[j - 1]
            i -= 1
            j -= 1

    Util.print_sequences(align1[::-1], align2[::-1])
    print("\tScore: " + str(score))


if __name__ == '__main__':
    seq_f = open("./tests/nucl_seqs/seq_nucl1").read()
    seq_s = open("./tests/nucl_seqs/seq_nucl2").read()
    calculate(seq_f, seq_s)
