package ru.bmstu.iu9.bioinf.commons;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class UtilsTest {

    @Test
    void testAminoAcidsScore() {
        assertEquals(1, Utils.aminoAcidsScore('N', 'D'));
        assertEquals(-3, Utils.aminoAcidsScore('I', 'R'));
    }

    @Test
    void testCodonPartialSubstitutionScore() {
        assertEquals(2, Utils.codonSubstitutionScore("AT_", 'L'));
        assertEquals(1.5, Utils.codonSubstitutionScore("_T_", 'M'));
        assertEquals(7, Utils.codonSubstitutionScore("TA_", 'Y'));
        assertEquals(1, Utils.codonSubstitutionScore("GTT", 'L'));
    }

}