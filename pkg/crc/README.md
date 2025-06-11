CRC即循环冗余校验码（Cyclic Redundancy Check）：是数据通信领域中最常用的一种查错校验码，其特征是信息字段和校验字段的长度可以任意选定。循环冗余检查（CRC）是一种数据传输检错功能，对数据进行多项式计算，并将得到的结果附在帧的后面，接收设备也执行类似的算法，以保证数据传输的正确性和完整性。

CRC算法参数模型解释:

```bash
NAME：参数模型名称。
WIDTH：宽度，即CRC比特数。
POLY：生成项的简写，以16进制表示。例如：CRC-32即是0x04C11DB7，忽略了最高位的"1"，即完整的生成项是0x104C11DB7。
INIT：这是算法开始时寄存器（crc）的初始化预置值，十六进制表示。
REFIN：待测数据的每个字节是否按位反转，True或False。
REFOUT：在计算后之后，异或输出之前，整个数据是否按位反转，True或False。
XOROUT：计算结果与此参数异或后得到最终的CRC值。
```

CRC参数模型

| CRC算法名称        | 多项式公式                                                   | 宽度 | 多项式     | 初始值     | 结果异或值 | 输入反转 | 输出反转 |
| ------------------ | ------------------------------------------------------------ | ---- | ---------- | ---------- | ---------- | -------- | -------- |
| CRC-4/ITU          | \(x^4 + x + 1\)                                              | 4    | 0x03       | 0x00       | 0x00       | true     | true     |
| CRC-5/EPC          | \(x^5 + x^3 + 1\)                                            | 5    | 0x09       | 0x09       | 0x00       | false    | false    |
| CRC-5/ITU          | \(x^5 + x^4 + x^2 + 1\)                                      | 5    | 0x15       | 0x00       | 0x00       | true     | true     |
| CRC-5/USB          | \(x^5 + x^2 + 1\)                                            | 5    | 0x05       | 0x1F       | 0x1F       | true     | true     |
| CRC-6/ITU          | \(x^6 + x + 1\)                                              | 6    | 0x03       | 0x00       | 0x00       | true     | true     |
| CRC-7/MMC          | \(x^7 + x^3 + 1\)                                            | 7    | 0x09       | 0x00       | 0x00       | false    | false    |
| CRC-8              | \(x^8 + x^2 + x + 1\)                                        | 8    | 0x07       | 0x00       | 0x00       | false    | false    |
| CRC-8/ITU          | \(x^8 + x^2 + x + 1\)                                        | 8    | 0x07       | 0x00       | 0x55       | false    | false    |
| CRC-8/ROHC         | \(x^8 + x^2 + x + 1\)                                        | 8    | 0x07       | 0xFF       | 0x00       | true     | true     |
| CRC-8/MAXIM        | \(x^8 + x^5 + x^4 + 1\)                                      | 8    | 0x31       | 0x00       | 0x00       | true     | true     |
| CRC-16/IBM         | \(x^{16} + x^{15} + x^2 + 1\)                                | 16   | 0x8005     | 0x0000     | 0x0000     | true     | true     |
| CRC-16/MAXIM       | \(x^{16} + x^{15} + x^2 + 1\)                                | 16   | 0x8005     | 0x0000     | 0xFFFF     | true     | true     |
| CRC-16/USB         | \(x^{16} + x^{15} + x^2 + 1\)                                | 16   | 0x8005     | 0xFFFF     | 0xFFFF     | true     | true     |
| CRC-16/MODBUS      | \(x^{16} + x^{15} + x^2 + 1\)                                | 16   | 0x8005     | 0xFFFF     | 0x0000     | true     | true     |
| CRC-16/CCITT       | \(x^{16} + x^{12} + x^5 + 1\)                                | 16   | 0x1021     | 0x0000     | 0x0000     | true     | true     |
| CRC-16/CCITT-FALSE | \(x^{16} + x^{12} + x^5 + 1\)                                | 16   | 0x1021     | 0xFFFF     | 0x0000     | false    | false    |
| CRC-16/X25         | \(x^{16} + x^{12} + x^5 + 1\)                                | 16   | 0x1021     | 0xFFFF     | 0xFFFF     | true     | true     |
| CRC-16/XMODEM      | \(x^{16} + x^{12} + x^5 + 1\)                                | 16   | 0x1021     | 0x0000     | 0x0000     | false    | false    |
| CRC-16/DNP         | \(x^{16} + x^{13} + x^{12} + x^{11} + x^{10} + x^8 + x^6 + x^5 + x^2 + 1\) | 16   | 0x3D65     | 0x0000     | 0xFFFF     | true     | true     |
| CRC-32             | \(x^{32} + x^{26} + x^{23} + x^{22} + x^{16} + x^{12} + x^{11} + x^{10} + x^8 + x^7 + x^5 + x^4 + x^2 + x + 1\) | 32   | 0x04C11DB7 | 0xFFFFFFFF | 0xFFFFFFFF | true     | true     |
| CRC-32/MPEG-2      | \(x^{32} + x^{26} + x^{23} + x^{22} + x^{16} + x^{12} + x^{11} + x^{10} + x^8 + x^7 + x^5 + x^4 + x^2 + x + 1\) | 32   | 0x04C11DB7 | 0xFFFFFFFF | 0x00000000 | false    | false    |