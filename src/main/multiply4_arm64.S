.text

.globl _vector_distance
_vector_distance:
    # REGISTERS
    # x0 => vector 1 (128 length vector)
    # x1 => vector 2
    # x3 => limit addr for x0
    # q0 => vector1 current lane
    # q1 => vector2 current lane
    # q2 => Accumulator of (error ^ 2)
    # Set the upper limit of vector
    mov x2, #32

    # zero out accumulator v2
    movi v2.4s, #0

    # iterate all vectors
    LOOP:
        # update q0, q1
        ldr q0, [x0], #16
        ldr q1, [x1], #16
        # sub q0 and q1 => q0 (error)
        sub v0.4s, v0.4s, v1.4s
        # q0 * q0 => q0 (error ^ 2)
        mul v0.4s, v0.4s, v0.4s
        # add (error ^ 2) to accumulator
        add v2.4s, v0.4s, v2.4s
        # check to loop again?
        subs x2, x2, #1
        b.ne LOOP
    
    # Unroll the accumulator
    umov w2, v2.s[0]
    umov w3, v2.s[1]
    umov w4, v2.s[2]
    umov w5, v2.s[3]
    # Add the unrolled values
    add w2, w2, w3
    add w4, w4, w5
    add w2, w2, w4
    sxtw x0, w2
    ret
    

.globl _vector_distance_8bit
_vector_distance_8bit:

    // x0 = vector1
    // x1 = vector2
    //
    // 128 int8 values
    // 16 int8 values processed per iteration
    // 8 iterations total
    //
    // v0  = vector1
    // v1  = vector2
    //
    // v2/v3 = lower 8 values widened to int16
    // v4/v5 = upper 8 values widened to int16
    //
    // v6-v9 = int32 differences
    // v10-v13 = squared differences
    //
    // v16/v17 = int64 accumulators

    movi v16.2d, #0
    movi v17.2d, #0

    // 128 / 16 = 8 iterations
    mov x2, #8

LOOP2:

    // ---------------------------------------------
    // Load 16 int8 values
    // ---------------------------------------------

    ldr q0, [x0], #16
    ldr q1, [x1], #16

    // ---------------------------------------------
    // Lower 8 int8 -> signed int16
    // ---------------------------------------------

    sshll v2.8h, v0.8b, #0
    sshll v3.8h, v1.8b, #0

    // v2 = vector1 - vector2
    sub v2.8h, v2.8h, v3.8h

    // ---------------------------------------------
    // Upper 8 int8 -> signed int16
    // ---------------------------------------------

    sshll2 v3.8h, v0.16b, #0
    sshll2 v4.8h, v1.16b, #0

    sub v3.8h, v3.8h, v4.8h

    // ---------------------------------------------
    // Widen lower differences:
    //
    // int16 -> int32
    // ---------------------------------------------

    sshll v6.4s, v2.4h, #0
    sshll2 v7.4s, v2.8h, #0

    // ---------------------------------------------
    // Widen upper differences:
    // ---------------------------------------------

    sshll v8.4s, v3.4h, #0
    sshll2 v9.4s, v3.8h, #0

    // ---------------------------------------------
    // Square the 32-bit differences
    // ---------------------------------------------

    mul v6.4s, v6.4s, v6.4s
    mul v7.4s, v7.4s, v7.4s
    mul v8.4s, v8.4s, v8.4s
    mul v9.4s, v9.4s, v9.4s

    // ---------------------------------------------
    // Horizontally add each group of 4 int32
    // ---------------------------------------------

    uaddlp v10.2d, v6.4s
    uaddlp v11.2d, v7.4s
    uaddlp v12.2d, v8.4s
    uaddlp v13.2d, v9.4s

    // ---------------------------------------------
    // Accumulate into 64-bit registers
    // ---------------------------------------------

    add v16.2d, v16.2d, v10.2d
    add v16.2d, v16.2d, v11.2d

    add v17.2d, v17.2d, v12.2d
    add v17.2d, v17.2d, v13.2d

    // ---------------------------------------------
    // Loop
    // ---------------------------------------------

    subs x2, x2, #1
    b.ne LOOP2

    // ---------------------------------------------
    // Combine both accumulators
    // ---------------------------------------------

    add v16.2d, v16.2d, v17.2d

    // v16.d[0] + v16.d[1]
    addp d16, v16.2d

    // Return int64 in x0
    umov x0, v16.d[0]

    ret