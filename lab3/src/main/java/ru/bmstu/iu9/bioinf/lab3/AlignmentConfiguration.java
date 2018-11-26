package ru.bmstu.iu9.bioinf.lab3;

class AlignmentConfiguration {
    private String nucleotideSeqFile, aminoAcidsFile;
    private String alignmentFile;
    private double open = -10;
    private double extend = -1;

    public void setAlignmentFile(String alignmentFile) {
        this.alignmentFile = alignmentFile;
    }


    public AlignmentConfiguration(String nucleotideSeqFile, String aminoAcidsFile, int open, int extend) {
        this.nucleotideSeqFile = nucleotideSeqFile;
        this.aminoAcidsFile = aminoAcidsFile;
//        this.alignmentFile = alignmentFile;
        this.open = open;
        this.extend = extend;
    }

    public AlignmentConfiguration(String nucleotideSeqFile, String aminoAcidsFile) {
        this.nucleotideSeqFile = nucleotideSeqFile;
        this.aminoAcidsFile = aminoAcidsFile;
    }


    public String getNucleotideSeqFile() {
        return nucleotideSeqFile;
    }

    public String getAminoAcidsFile() {
        return aminoAcidsFile;
    }

    public String getAlignmentFile() {
        return alignmentFile;
    }

    public double getExtend() {
        return extend;
    }

    public double getOpen() {
        return open;
    }

    public void setNucleotideSeqFile(String nucleotideSeqFile) {
        this.nucleotideSeqFile = nucleotideSeqFile;
    }

    public void setAminoAcidsFile(String aminoAcidsFile) {
        this.aminoAcidsFile = aminoAcidsFile;
    }

    public void setOpen(double open) {
        this.open = open;
    }

    public void setExtend(double extend) {
        this.extend = extend;
    }
}