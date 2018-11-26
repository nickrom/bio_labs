package ru.bmstu.iu9.bioinf.commons;

import java.io.IOException;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashSet;
import java.util.Set;

public class Utils {
    public static int aminoAcidsScore(char a, char b) {
        int aIdx = Consts.AMINO_ACID_2_IDX.get(a);
        int bIdx = Consts.AMINO_ACID_2_IDX.get(b);

        return Consts.BLOSUM62_MATRIX[aIdx][bIdx];
    }

    private static int codonFullSubstitutionScore(String codon, char nucleotide) {
        int curNuclIdx = Consts.AMINO_ACID_2_IDX.get(nucleotide);
        int otherNuclIdx = Consts.AMINO_ACID_2_IDX.get(Consts.CODON_2_AMINO_ACID.get(codon));

        return Consts.BLOSUM62_MATRIX[curNuclIdx][otherNuclIdx];
    }

    public static double codonSubstitutionScore(String codonPtr, char aminoAcid) {
        if(!codonPtr.contains("_")) {
            return codonFullSubstitutionScore(codonPtr, aminoAcid);
        }

        Set<String> codonsSubtitutionsSet = new HashSet<>();
        codonsSubtitutionsSet.add(codonPtr);
        int prevSetSize = 0;

        while(prevSetSize != codonsSubtitutionsSet.size()) {
            prevSetSize = codonsSubtitutionsSet.size();
            Set<String> newGeneration = new HashSet<>();
            for(String c : codonsSubtitutionsSet) {
                int slashIdx = c.indexOf('_');
                if(slashIdx >= 0) {
                    for (int i = 0; i < Consts.NUCLEOTIDES.length; i++) {
                        StringBuilder newCodon = new StringBuilder();

                        newCodon.append(c, 0, slashIdx);
                        newCodon.append(Consts.NUCLEOTIDES[i]);
                        newCodon.append(c, slashIdx + 1, c.length());

                        newGeneration.add(newCodon.toString());
                    }
                }
            }
            codonsSubtitutionsSet.addAll(newGeneration);
        }

        int totalScore = 0, i = 0;

        for(String codon : codonsSubtitutionsSet) {
            if(!codon.contains("_") && Consts.CODON_2_AMINO_ACID.get(codon) != '#') {
                totalScore += codonFullSubstitutionScore(codon, aminoAcid);
                i++;
            }
        }

        return i == 0 ? 0 : (double) totalScore / (double) i;
    }

    public static String readSequence(String file) throws IOException {
        return new String(Files.readAllBytes(Paths.get(file)));
    }

    public static void printMatrix(PrintWriter out, String matrixName, double[][] matrix, String seq1, String seq2) {
        int m = matrix[0].length;
        int n = matrix.length;

        out.println("\t MATRIX " + (matrixName == null ? "" : "'" + matrixName + "'") + ":\n");

        out.print("#    ");
        for (int i = 0; i < m; i++) {
            out.print((i > 0 ? seq2.charAt(i - 1) : '_') + " | ");
        }
        out.println("\n" + new String(new char[4*(m + 1) - 1]).replace("\0", "_"));
        for (int i = 0; i < n; i++) {
            out.print((i > 0 ? seq1.charAt(i - 1) : '_') + " | ");
            for (int j = 0; j < m; j++) {
                out.printf("%f\t", matrix[i][j]);
            }
            out.println();
        }

        out.println(new String(new char[4*(m + 1) - 1]).replace("\0", "_"));
    }
}
