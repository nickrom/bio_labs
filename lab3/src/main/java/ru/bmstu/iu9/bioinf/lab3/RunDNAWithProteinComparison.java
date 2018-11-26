package ru.bmstu.iu9.bioinf.lab3;

import org.apache.commons.cli.*;

import javax.naming.ConfigurationException;
import java.io.File;
import java.io.IOException;
import java.io.PrintWriter;
import java.util.Optional;

import static ru.bmstu.iu9.bioinf.commons.Utils.codonSubstitutionScore;
import static ru.bmstu.iu9.bioinf.commons.Utils.printMatrix;
import static ru.bmstu.iu9.bioinf.commons.Utils.readSequence;

public class RunDNAWithProteinComparison {
    private static CommandLine cmd;
    private static final double MIN_INF = Double.MIN_VALUE;
    private static final PrintWriter SYS_OUT_WRITER = new PrintWriter(System.out, true);

    public static void main(String[] args) {
        try {
            initCmdOpts(args);
            AlignmentConfiguration conf = getAlignConfigFromCmd();
            alignSequences(conf);
        } catch (Exception e) {
            if(!(e instanceof ParseException)) {
                e.printStackTrace();
            }
        }
    }

    private static void alignSequences(AlignmentConfiguration conf) throws IOException {
        String aminoAcidsSeq = readSequence(conf.getAminoAcidsFile());
        String nucleotideSeq = readSequence(conf.getNucleotideSeqFile());

//        String alignmentFile = conf.getAlignmentFile();
//        PrintWriter output = alignmentFile == null ? SYS_OUT_WRITER : new PrintWriter(alignmentFile);

        int n = nucleotideSeq.length() + 1, m = aminoAcidsSeq.length() + 1;
        double gap = conf.getOpen();

        double[][] delNucMatrix = new double[n][m];
        double[][] insAaMatrix = new double[n][m];
        double[][] substMatrix = new double[n][m];

        // init matrices
        substMatrix[0][0] = 0;
        delNucMatrix[0][0] = insAaMatrix[0][0] = MIN_INF;
        for (int i = 1; i < n; i++) {
            substMatrix[i][0] = MIN_INF;
            delNucMatrix[i][0] = insAaMatrix[i][0] = MIN_INF;
        }

        for (int j = 0; j < m; j++) {
            insAaMatrix[0][j] = MIN_INF;
            delNucMatrix[0][j] = MIN_INF;
        }

        System.out.println("Before alignment: ");
        printMatrix(SYS_OUT_WRITER, "Substitution", substMatrix, nucleotideSeq, aminoAcidsSeq);
        printMatrix(SYS_OUT_WRITER, "Nucleotide deletion", delNucMatrix, nucleotideSeq, aminoAcidsSeq);
        printMatrix(SYS_OUT_WRITER, "Amino acids insertion", insAaMatrix, nucleotideSeq, aminoAcidsSeq);


        for (int j = 1; j < m; j++) {
            for (int i = 1; i < n - 3; i++) {
                String codon = substitutionAt(nucleotideSeq, i - 1);
                char aminoAcid = aminoAcidsSeq.charAt(j - 1);

                double delScore = delScore(nucleotideSeq, i - 1, aminoAcid);
                double subScore = codonSubstitutionScore(codon, aminoAcid);

                delNucMatrix[i][j] = max(
                        delNucMatrix[i - 1][j] + delScore,
                        substMatrix[i - 1][j] + delScore,
                        insAaMatrix[i - 1][j] + delScore
                );
                insAaMatrix[i][j] = max(
                        delNucMatrix[i][j - 1] + gap,
                        substMatrix[i][j - 1] + gap,
                        insAaMatrix[i][j - 1] + gap
                );
                substMatrix[i][j] = max(
                        delNucMatrix[i - 1][j - 1] + subScore,
                        substMatrix[i - 1][j - 1] + subScore,
                        insAaMatrix[i - 1][j - 1] + subScore
                );
            }
        }

        // TODO: 27.11.18 Remove debug prints
        System.out.println("After alignment: ");
        printMatrix(SYS_OUT_WRITER, "Substitution", substMatrix, nucleotideSeq, aminoAcidsSeq);
        printMatrix(SYS_OUT_WRITER, "Nucleotide deletion", delNucMatrix, nucleotideSeq, aminoAcidsSeq);
        printMatrix(SYS_OUT_WRITER, "Amino acids insertion", insAaMatrix, nucleotideSeq, aminoAcidsSeq);
    }

    private static String substitutionAt(String nucleotideSeq, int i) {
        int n = nucleotideSeq.length();
        StringBuilder codon = new StringBuilder();

        codon.append(i < n ? nucleotideSeq.charAt(i) : "_");
        codon.append(i + 1 < n ? nucleotideSeq.charAt(i + 1) : "_");
        codon.append(i + 2 < n ? nucleotideSeq.charAt(i + 2) : "_");

        return codon.toString();
    }

    private static double delScore(String nucleotideSeq, int i, char aminoAcid) {
        int n = nucleotideSeq.length();
        StringBuilder codon = new StringBuilder();

        codon.append("_");
        codon.append(i + 1 < n ? nucleotideSeq.charAt(i + 1) : "_");
        codon.append(i + 2 < n ? nucleotideSeq.charAt(i + 2) : "_");

        return codonSubstitutionScore(codon.toString(), aminoAcid);
    }

    private static double max(double ...nums) {
        double ret = Double.MIN_VALUE;
        for(double num : nums) {
            if(num > ret)
                ret = num;
        }

        return ret;
    }

    private static void initCmdOpts(String[] args) throws ParseException {
        Options cmdOptions = new Options();

        cmdOptions.addOption(
                Option.builder("h")
                        .longOpt("help")
                        .desc("Print help message")
                        .hasArg(false)
                        .build()
        );
        cmdOptions.addOption(
                Option.builder("n")
                        .longOpt("nucleotide-sequence")
                        .desc("Input file with nucleotide sequence")
                        .hasArgs()
                        .numberOfArgs(1)
                        .required()
                        .build()
        );
        cmdOptions.addOption(
                Option.builder("a")
                        .longOpt("amino-acid-sequence")
                        .desc("Input file with amino acids' sequence")
                        .hasArgs()
                        .numberOfArgs(1)
                        .required()
                        .build()
        );


        cmdOptions.addOption(
                Option.builder("g")
                        .longOpt("gap")
                        .desc("Gap penalty.")
                        .hasArg()
                        .numberOfArgs(1)
                        .type(Integer.class)
                        .required()
                        .build()
        );

        cmdOptions.addOption(
                Option.builder("o")
                        .desc("Output file with result alignment and score.")
                        .longOpt("output")
                        .hasArg()
                        .type(String.class)
                        .build()
        );

        try {
            cmd = new DefaultParser().parse(cmdOptions, args);
        } catch (ParseException e) {
            String msg = e.getMessage();
            System.err.printf("[error] Invalid command line arguments%s.\n", (msg == null || msg.length() == 0) ? "" : ": " + msg);
            HelpFormatter formatter = new HelpFormatter();
            String executableName = new File(
                    RunDNAWithProteinComparison.class.getProtectionDomain()
                            .getCodeSource()
                            .getLocation()
                            .getPath()
            )
                    .getName();
            formatter.printHelp(executableName, cmdOptions, true);

            throw e;
        }
    }

    private static AlignmentConfiguration getAlignConfigFromCmd() throws ConfigurationException {
        String nucleotideSeqFile, aminoAcidSeqFile;
        int gap;

        if(!cmd.hasOption("a") || !cmd.hasOption("n")) {
            throw new ConfigurationException("Invalid number of input files.");
        }

        nucleotideSeqFile = cmd.getOptionValue("n");
        aminoAcidSeqFile = cmd.getOptionValue("a");

        gap = parseInt(cmd.getOptionValue("gap")).orElse(-10);

        if(gap > 0) {
            throw new ConfigurationException("Invalid open value. It should be negative integer");
        }

        AlignmentConfiguration conf = new AlignmentConfiguration(nucleotideSeqFile, aminoAcidSeqFile);
        conf.setOpen(gap);

        if(cmd.hasOption('o')) {
            conf.setAlignmentFile(cmd.getOptionValue('o'));
        }

        return conf;
    }

    private static Optional<Integer> parseInt(String toParse) {
        try {
            return Optional.of(Integer.parseInt(toParse));
        } catch (NumberFormatException e) {
            return Optional.empty();
        }
    }
}
