TEXT ·SquaredDistanceNEON(SB), $0-56
    // R0 --> query ptr
    // R1 --> vector ptr
    // R2 --> Length/Limit of query ptr
    // V0 --> query lane (4x int32)
    // V1 --> vector lane (4x int32)
    // V2 --> diff lane
    // R3 --> result (int64 accumulator — avoids overflow on real data)
    // R4 + onwards calculations

    MOVD query+0(FP), R0
    MOVD vector+24(FP), R1
    MOVD query+8(FP), R2
    LSL $2, R2, R2
    ADD R0, R2, R2
    MOVD $0, R3

    Loop:
        VLD1.P 16(R0), [V0.S4]
        VLD1.P 16(R1), [V1.S4]
        VSUB V0.S4, V1.S4, V2.S4

        // Lane 0
        VMOV V2.S[0], R4
        MULW R4, R4, R4

        // Lane 1
        VMOV V2.S[1], R5
        MULW R5, R5, R5

        // Lane 2
        VMOV V2.S[2], R6
        MULW R6, R6, R6

        // Lane 3
        VMOV V2.S[3], R7
        MULW R7, R7, R7

        // Now combine the independent results
        ADD R4, R5, R4
        ADD R6, R7, R6
        ADD R4, R6, R4

        ADD R4, R3, R3

        CMP R2, R0
        BLO Loop

    MOVD R3, ret+48(FP)   // full 64-bit store now
    RET
