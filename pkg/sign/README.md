# AES配套的其它语种调用

### java

```java
import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.SecureRandom;
import java.util.Arrays;
import java.util.Base64;

/**
 * AES 加密工具类
 */
public class AesUtil {

    // AES 算法
    private static final String AES = "AES";
    // AES CBC 模式，PKCS5Padding 填充
    private static final String AES_CBC_PKCS5PADDING = "AES/CBC/PKCS5Padding";
    // SHA-256 消息摘要算法
    private static final String SHA3_256 = "SHA-256";

    /**
     * 根据密码生成密钥
     * 
     * @param password 密码
     * @param length   密钥长度
     * @return 密钥
     * @throws Exception 如果生成密钥失败
     */
    public static byte[] generateKey(String password, int length) throws Exception {
        // 创建消息摘要对象
        MessageDigest digest = MessageDigest.getInstance(SHA3_256);
        // 计算密码的哈希值
        byte[] key = digest.digest(password.getBytes(StandardCharsets.UTF_8));
        // 截取哈希值的前 length 字节作为密钥
        return Arrays.copyOf(key, length);
    }

    /**
     * 使用 AES 算法加密文本
     * 
     * @param plainText 明文
     * @param key       密钥
     * @return 密文
     * @throws Exception 如果加密失败
     */
    public static String encrypt(String plainText, byte[] key) throws Exception {
        // 检查密钥是否为空
        if (key.length == 0) {
            throw new Exception("密钥不能为空");
        }

        // 创建 Cipher 对象
        Cipher cipher = Cipher.getInstance(AES_CBC_PKCS5PADDING);
        // 创建密钥规格对象
        SecretKeySpec secretKey = new SecretKeySpec(key, AES);

        // 生成随机初始化向量
        byte[] iv = new byte[cipher.getBlockSize()];
        new SecureRandom().nextBytes(iv);

        // 初始化 Cipher 对象为加密模式
        cipher.init(Cipher.ENCRYPT_MODE, secretKey, new IvParameterSpec(iv));

        // 将明文转换为字节数组
        byte[] plainTextBytes = plainText.getBytes(StandardCharsets.UTF_8);

        // 加密字节数组
        byte[] encrypted = cipher.doFinal(plainTextBytes);

        // 将初始化向量和加密结果合并
        byte[] cipherTextWithIv = new byte[iv.length + encrypted.length];
        System.arraycopy(iv, 0, cipherTextWithIv, 0, iv.length);
        System.arraycopy(encrypted, 0, cipherTextWithIv, iv.length, encrypted.length);

        // 将合并结果进行 Base64 编码
        return Base64.getEncoder().encodeToString(cipherTextWithIv);
    }

    /**
     * 使用 AES 算法解密文本
     * 
     * @param cipherText 密文
     * @param key        密钥
     * @return 明文
     * @throws Exception 如果解密失败
     */
    public static String decrypt(String cipherText, byte[] key) throws Exception {
        // 检查密钥是否为空
        if (key.length == 0) {
            throw new Exception("密钥不能为空");
        }

        // 将密文进行 Base64 解码
        byte[] cipherTextBytes = Base64.getDecoder().decode(cipherText);

        // 创建 Cipher 对象
        Cipher cipher = Cipher.getInstance(AES_CBC_PKCS5PADDING);
        // 创建密钥规格对象
        SecretKeySpec secretKey = new SecretKeySpec(key, AES);

        // 获取 Cipher 对象的块大小
        int blockSize = cipher.getBlockSize();
        // 提取初始化向量
        byte[] iv = Arrays.copyOfRange(cipherTextBytes, 0, blockSize);

        // 初始化 Cipher 对象为解密模式
        cipher.init(Cipher.DECRYPT_MODE, secretKey, new IvParameterSpec(iv));

        // 解密字节数组
        byte[] original = cipher.doFinal(Arrays.copyOfRange(cipherTextBytes, blockSize, cipherTextBytes.length));

        // 将解密结果转换为字符串
        return new String(original, StandardCharsets.UTF_8);
    }
}

// 测试类如下所示
public class AesUtilTest {

    public static void main(String[] args) throws Exception {
        String password = "mysecretpassword";
        int keyLength = 16; // AES-128
        byte[] key = AesUtil.generateKey(password, keyLength);

        Object[] originalTexts = {
            "Hello, World!", // 字符串
            12345, // 整数
            3.14159265359, // 浮点数
            true, // 布尔值
            new byte[] {1, 2, 3, 4, 5}, // 字节数组
            new int[] {1, 2, 3, 4, 5}, // 整型数组
            new double[] {1.1, 2.2, 3.3, 4.4, 5.5}, // 浮点数数组
			"中文测试",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
        };

        for (Object originalText : originalTexts) {
            System.out.println("Original Text: " + originalText);

            // 转换为字符串
            String originalTextStr = String.valueOf(originalText);

            // 加密
            String encryptedText = AesUtil.encrypt(originalTextStr, key);
            System.out.println("Encrypted Text: " + encryptedText);

            // 解密
            String decryptedText = AesUtil.decrypt(encryptedText, key);
            System.out.println("Decrypted Text: " + decryptedText);

            // 验证
            if (originalTextStr.equals(decryptedText)) {
                System.out.println("Success: Decrypted text matches the original text.");
            } else {
                System.out.println("Error: Decrypted text does not match the original text.");
            }

            System.out.println();
        }
    }
}
```