

TEXT ·TestAssemblyParameter(SB), $0-32
    // R0 --> PTR of ARRAY (CUR)
    // R1 --> LEN of ARRAY/ LIMIT of Loop = ARRAY LEN * 8
    // V0 --> Lane 0 Cur Element
    // V1 --> Lane 1 Cur Element
    // V2 --> Lane 0 Cur Element
    // V3 --> Lane 1 Cur Element
    // R2 --> Result of Lane 0
    // R3 --> Result of Lane 1
    // R4 --> final result

    // Read the pointer for array --> ptr
    MOVD a+0(FP), R0

    // Create the limit of the Loop
    MOVD a+8(FP), R1
    LSL $3, R1, R1
    ADD R1, R0, R1

    // Create Running Sums
    VEOR V2.B16, V2.B16, V2.B16
    VEOR V3.B16, V3.B16, V3.B16

    // Loop all elements
    Loop:
        // Extract cur Array
        VLD1.P 32(R0), [V0.D2, V1.D2]
        VADD V0.D2, V2.D2, V2.D2
        VADD V1.D2, V3.D2, V3.D2

        CMP R1, R0
        BLO Loop
    
    VMOV V2.D[0], R2
    VMOV V2.D[1], R3
    ADD R2, R3, R4

    VMOV V3.D[0], R2
    VMOV V3.D[1], R3
    ADD R2, R3, R2
    ADD R2, R4, R4

    // Return Sum
    MOVD R4, ret+24(FP)

    RET
