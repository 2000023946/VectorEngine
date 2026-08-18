TEXT ·SquaredDistanceNEON(SB), $0-56
    // R0 --> query ptr
    // R1 --> vector ptr
    // R2 --> end pointer
    // R3 --> final result
    //
    // V0 --> query values 0-3
    // V1 --> query values 4-7
    // V2 --> vector values 0-3
    // V3 --> vector values 4-7
    // V4 --> differences 0-3
    // V5 --> differences 4-7
    //
    // R4-R7 --> scalar squared lane results

    MOVD query+0(FP), R0
    MOVD vector+24(FP), R1

    // Calculate end pointer:
    // end = query_ptr + length * 4
    MOVD query+8(FP), R2
    LSL $2, R2, R2
    ADD R0, R2, R2

    // Final accumulator
    MOVD $0, R3

Loop:
    // Load 8 int32 values from query
    VLD1.P 32(R0), [V0.S4, V1.S4]

    // Load 8 int32 values from vector
    VLD1.P 32(R1), [V2.S4, V3.S4]

    // Calculate differences for both groups
    VSUB V0.S4, V2.S4, V4.S4
    VSUB V1.S4, V3.S4, V5.S4

    // -------------------------
    // First 4 lanes
    // -------------------------

    VMOV V4.S[0], R4
    MULW R4, R4, R4

    VMOV V4.S[1], R5
    MULW R5, R5, R5

    VMOV V4.S[2], R6
    MULW R6, R6, R6

    VMOV V4.S[3], R7
    MULW R7, R7, R7

    ADD R4, R5, R4
    ADD R6, R7, R6
    ADD R4, R6, R4

    ADD R4, R3, R3

    // -------------------------
    // Second 4 lanes
    // -------------------------

    VMOV V5.S[0], R4
    MULW R4, R4, R4

    VMOV V5.S[1], R5
    MULW R5, R5, R5

    VMOV V5.S[2], R6
    MULW R6, R6, R6

    VMOV V5.S[3], R7
    MULW R7, R7, R7

    ADD R4, R5, R4
    ADD R6, R7, R6
    ADD R4, R6, R4

    ADD R4, R3, R3

    CMP R2, R0
    BLO Loop

    MOVD R3, ret+48(FP)
    RET
