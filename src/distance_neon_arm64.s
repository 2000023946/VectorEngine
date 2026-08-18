TEXT ·SquaredDistanceNEON(SB), $0-56
    // R0  --> query ptr
    // R1  --> vector ptr
    // R2  --> end pointer
    // R3  --> final int64 accumulator
    //
    // V0 --> query values 0-7 (int8)
    // V1 --> vector values 0-7 (int8)
    //
    // Each lane gets independent registers so the CPU
    // can potentially execute the arithmetic in parallel.

    MOVD query+0(FP), R0
    MOVD vector+24(FP), R1

    // int8 = 1 byte
    //
    // end = query_ptr + length
    MOVD query+8(FP), R2
    ADD R0, R2, R2

    // Final accumulator
    MOVD $0, R3

Loop:
    // Load 8 int8 values from query
    VLD1.P 8(R0), [V0.B8]

    // Load 8 int8 values from vector
    VLD1.P 8(R1), [V1.B8]

    // --------------------------------------------------
    // Lane 0
    // --------------------------------------------------

    VMOV V0.B[0], R4
    VMOV V1.B[0], R5
    SXTB R4, R4
    SXTB R5, R5
    SUB R5, R4, R4
    MULW R4, R4, R4

    // --------------------------------------------------
    // Lane 1
    // --------------------------------------------------

    VMOV V0.B[1], R6
    VMOV V1.B[1], R7
    SXTB R6, R6
    SXTB R7, R7
    SUB R7, R6, R6
    MULW R6, R6, R6

    // --------------------------------------------------
    // Lane 2
    // --------------------------------------------------

    VMOV V0.B[2], R8
    VMOV V1.B[2], R9
    SXTB R8, R8
    SXTB R9, R9
    SUB R9, R8, R8
    MULW R8, R8, R8

    // --------------------------------------------------
    // Lane 3
    // --------------------------------------------------

    VMOV V0.B[3], R10
    VMOV V1.B[3], R11
    SXTB R10, R10
    SXTB R11, R11
    SUB R11, R10, R10
    MULW R10, R10, R10

    // --------------------------------------------------
    // Reduce lanes 0-3
    // --------------------------------------------------

    ADD R6, R4, R4
    ADD R8, R4, R4
    ADD R10, R4, R4
    ADD R4, R3, R3

    // --------------------------------------------------
    // Lane 4
    // --------------------------------------------------

    VMOV V0.B[4], R4
    VMOV V1.B[4], R5
    SXTB R4, R4
    SXTB R5, R5
    SUB R5, R4, R4
    MULW R4, R4, R4

    // --------------------------------------------------
    // Lane 5
    // --------------------------------------------------

    VMOV V0.B[5], R6
    VMOV V1.B[5], R7
    SXTB R6, R6
    SXTB R7, R7
    SUB R7, R6, R6
    MULW R6, R6, R6

    // --------------------------------------------------
    // Lane 6
    // --------------------------------------------------

    VMOV V0.B[6], R8
    VMOV V1.B[6], R9
    SXTB R8, R8
    SXTB R9, R9
    SUB R9, R8, R8
    MULW R8, R8, R8

    // --------------------------------------------------
    // Lane 7
    // --------------------------------------------------

    VMOV V0.B[7], R10
    VMOV V1.B[7], R11
    SXTB R10, R10
    SXTB R11, R11
    SUB R11, R10, R10
    MULW R10, R10, R10

    // --------------------------------------------------
    // Reduce lanes 4-7
    // --------------------------------------------------

    ADD R6, R4, R4
    ADD R8, R4, R4
    ADD R10, R4, R4
    ADD R4, R3, R3

    // --------------------------------------------------
    // Continue until all dimensions are processed.
    // --------------------------------------------------

    CMP R2, R0
    BLO Loop

    MOVD R3, ret+48(FP)
    RET
