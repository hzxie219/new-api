# Golang 安全编码 Checklist

本checklist基于golang安全编码规范整理，用于代码编写时的安全自查和安全代码审查。

## 使用说明

### 优先级说明
- **【强制】**: 必须遵守，违反会导致严重的安全风险
- **【可选】**: 建议遵守，有助于提高代码安全性
- **【建议】**: 推荐遵守，有助于提升整体安全质量

### 使用场景
1. **编码前**: 了解相关安全规范要求
2. **编码中**: 实时参考对应的安全条目
3. **编码后**: 逐项检查所有【强制】标记的条目
4. **安全审查**: 使用本checklist进行系统性安全审查
5. **持续改进**: 定期回顾，养成安全编码习惯

---

## 目录
- [01. 内存管理](#01-内存管理)
- [02. 文件操作](#02-文件操作)
- [03. 系统接口](#03-系统接口)
- [04. 通讯安全](#04-通讯安全)
- [05. 敏感数据保护](#05-敏感数据保护)
- [06. 加密解密](#06-加密解密)
- [07. 正则表达式](#07-正则表达式)
- [08. 输入输出校验](#08-输入输出校验)
- [09. SQL操作](#09-sql操作)
- [10. 网络请求](#10-网络请求)
- [11. 服务器端渲染](#11-服务器端渲染)
- [12. Web跨域](#12-web跨域)
- [13. 响应输出](#13-响应输出)
- [14. Session安全](#14-session安全)
- [15. 访问控制](#15-访问控制)
- [16. 并发保护](#16-并发保护)
- [17. 异常行为](#17-异常行为)
- [18. 序列化与反序列化](#18-序列化与反序列化)

---

## 01. 内存管理

### 1.1 切片长度校验 【强制】
- [ ] 对slice进行操作前，必须判断长度是否合法，防止程序panic
- [ ] 使用slice索引访问前检查len(slice)
- [ ] 使用slice切片操作前验证边界合法性

```go
// ❌ 错误示例
func decode(data []byte) bool {
    if data[0] == 'F' && data[1] == 'U' {  // 未检查长度
        return true
    }
    return false
}

// ✅ 正确示例
func decode(data []byte) bool {
    if len(data) >= 2 {
        if data[0] == 'F' && data[1] == 'U' {
            return true
        }
    }
    return false
}
```

### 1.2 nil指针判断 【强制】
- [ ] 进行指针操作时，必须判断该指针是否为nil
- [ ] 结构体Unmarshal后，访问指针字段前必须检查nil
- [ ] 函数返回指针类型时，调用方必须进行nil检查

```go
// ❌ 错误示例
func main() {
    packet := new(Packet)
    packet.UnmarshalBinary(data)
    fmt.Printf("Stat: %v\n", packet.Data.Stat)  // Data可能为nil
}

// ✅ 正确示例
func main() {
    packet := new(Packet)
    packet.UnmarshalBinary(data)
    if packet.Data != nil {
        fmt.Printf("Stat: %v\n", packet.Data.Stat)
    }
}
```

### 1.3 整数安全 【强制】

#### 1.3.1 确保无符号整数运算时不会反转
- [ ] 无符号整数运算前检查是否会发生反转
- [ ] 涉及数组索引、对象长度、数组边界时严格校验
- [ ] 使用math.MaxUint/MinUint常量进行边界检查

```go
// ❌ 错误示例
func sum(a, b uint64) uint64 {
    return a + b  // 可能反转
}

// ✅ 正确示例
func sum(a, b uint64) (uint64, error) {
    if math.MaxUint64-a < b {
        return 0, errors.New("unsigned integer overflow")
    }
    return a + b, nil
}
```

#### 1.3.2 确保有符号整数运算时不会溢出
- [ ] 有符号整数运算前检查是否会溢出
- [ ] 涉及数组索引、对象长度、数组边界时严格校验
- [ ] 使用math.MaxInt/MinInt常量进行边界检查

```go
// ❌ 错误示例
func add(a, b int32) int32 {
    return a + b  // 可能溢出
}

// ✅ 正确示例
func add(a, b int32) (int32, error) {
    if (a > 0 && b > math.MaxInt32-a) || (b < 0 && a < math.MinInt32-b) {
        return 0, errors.New("integer overflow")
    }
    return a + b, nil
}
```

#### 1.3.3 确保整型转换时不会出现截断错误
- [ ] 较大整型转较小整型时，校验数据范围
- [ ] 使用math包常量验证转换是否安全
- [ ] 涉及数组索引、长度、边界时必须严格校验

```go
// ❌ 错误示例
func convert(a int32) int16 {
    return int16(a)  // 可能截断
}

// ✅ 正确示例
func convert(a int32) (int16, error) {
    if a < math.MinInt16 || a > math.MaxInt16 {
        return 0, errors.New("integer truncation")
    }
    return int16(a), nil
}
```

#### 1.3.4 确保整型转换时不会出现符号错误
- [ ] 有符号整型转无符号整型前检查是否为负数
- [ ] 无符号整型转有符号整型前检查是否超出范围

```go
// ❌ 错误示例
func convert(a int32) uint32 {
    return uint32(a)  // a为负数时会变成大正数
}

// ✅ 正确示例
func convert(a int32) (uint32, error) {
    if a < 0 {
        return 0, errors.New("negative value cannot convert to unsigned")
    }
    return uint32(a), nil
}
```

### 1.4 make分配长度验证 【强制】
- [ ] 使用make创建slice/map时，对外部输入的长度进行合法性校验
- [ ] 检查长度是否为负数
- [ ] 检查长度是否超过合理上限（如64MB）
- [ ] 防止外部输入导致程序panic或内存耗尽

```go
// ❌ 错误示例
func parse(size int, data []byte) []byte {
    buffer := make([]byte, size)  // size未校验
    copy(buffer, data)
    return buffer
}

// ✅ 正确示例
func parse(size int, data []byte) ([]byte, error) {
    if size < 0 || size > 64*1024*1024 {
        return nil, errors.New("invalid size")
    }
    buffer := make([]byte, size)
    copy(buffer, data)
    return buffer, nil
}
```

### 1.5 禁止SetFinalizer和指针循环引用同时使用 【强制】
- [ ] 避免循环引用与runtime.SetFinalizer结合使用
- [ ] 使用SetFinalizer时确保对象无循环引用
- [ ] 注意内存泄漏风险

### 1.6 禁止重复释放channel 【强制】
- [ ] 避免在异常流程中重复关闭channel
- [ ] 使用defer关闭channel确保只关闭一次
- [ ] 使用sync.Once确保channel只关闭一次

```go
// ❌ 错误示例
func foo(c chan int) {
    defer close(c)
    if err := processBusiness(); err != nil {
        close(c)  // 重复关闭
        return
    }
}

// ✅ 正确示例
func foo(c chan int) {
    defer close(c)  // 只关闭一次
    if err := processBusiness(); err != nil {
        return
    }
}
```

### 1.7 确保对channel是否关闭做检查 【可选】
- [ ] 从channel读取数据时检查是否已关闭
- [ ] 使用comma-ok模式检查channel状态
- [ ] 避免因channel关闭导致死循环

```go
// ❌ 错误示例
for {
    select {
    case <-cc:  // channel关闭后会一直触发
        fmt.Println("continue")
    }
}

// ✅ 正确示例
for {
    select {
    case _, ok := <-cc:
        if !ok {
            return  // channel已关闭
        }
        fmt.Println("continue")
    }
}
```

### 1.8 确保每个协程都能退出 【强制】
- [ ] 每个goroutine都必须有退出条件
- [ ] 使用context或channel控制goroutine生命周期
- [ ] 使用sync.WaitGroup等待goroutine完成
- [ ] 防止goroutine泄漏导致内存泄漏

```go
// ❌ 错误示例
func worker(name string) {
    for {  // 永远不会退出
        time.Sleep(1 * time.Second)
        fmt.Println(name)
    }
}

// ✅ 正确示例
func worker(name string, ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return  // 收到退出信号
        default:
            time.Sleep(1 * time.Second)
            fmt.Println(name)
        }
    }
}
```

### 1.9 不使用unsafe包 【可选】
- [ ] 避免使用unsafe包绕过Go内存安全机制
- [ ] 必须使用时做好完整的安全校验
- [ ] 记录使用unsafe包的理由和风险评估

### 1.10 不使用slice作为函数入参 【强制】
- [ ] 慎用slice作为函数入参（slice是引用类型）
- [ ] 明确函数是否会修改入参slice
- [ ] 必要时使用copy创建slice副本

---

## 02. 文件操作

### 2.1 确保文件路径验证前对其进行标准化 【强制】
- [ ] 使用filepath.Clean标准化文件路径
- [ ] 验证路径不存在目录穿越（../）
- [ ] 使用filepath.Abs获取绝对路径
- [ ] 确保路径在预期的目录范围内

```go
// ✅ 正确示例
func validatePath(userPath string) (string, error) {
    cleanPath := filepath.Clean(userPath)
    absPath, err := filepath.Abs(cleanPath)
    if err != nil {
        return "", err
    }

    // 检查是否在允许的目录内
    if !strings.HasPrefix(absPath, allowedDir) {
        return "", errors.New("path traversal detected")
    }
    return absPath, nil
}
```

### 2.2 确保在多用户系统中创建文件时指定合适的访问许可 【强制】
- [ ] 创建文件时显式指定文件权限
- [ ] 敏感文件权限设置为0600（仅所有者可读写）
- [ ] 避免使用过于宽松的权限（如0777）
- [ ] 使用os.OpenFile指定权限模式

```go
// ✅ 正确示例
file, err := os.OpenFile("sensitive.dat", os.O_CREATE|os.O_WRONLY, 0600)
```

### 2.3 避免在共享目录操作文件 【强制】
- [ ] 避免在/tmp等共享目录创建可预测名称的文件
- [ ] 使用os.CreateTemp创建临时文件
- [ ] 防止符号链接攻击和竞态条件
- [ ] 检查文件是否已存在

### 2.4 确保安全地从压缩包中提取文件 【强制】
- [ ] 解压前验证文件路径，防止目录穿越
- [ ] 检查解压后的文件路径是否在预期目录内
- [ ] 限制解压文件的大小，防止zip炸弹
- [ ] 使用filepath.Clean清理路径

```go
// ✅ 正确示例
func extractFile(f *zip.File, destDir string) error {
    filePath := filepath.Join(destDir, f.Name)
    cleanPath := filepath.Clean(filePath)

    if !strings.HasPrefix(cleanPath, filepath.Clean(destDir)) {
        return errors.New("illegal file path")
    }

    // 限制文件大小
    if f.UncompressedSize64 > maxFileSize {
        return errors.New("file too large")
    }

    // 继续解压...
    return nil
}
```

### 2.5 确保临时文件使用完毕后及时删除 【强制】
- [ ] 使用defer确保临时文件被删除
- [ ] 使用os.CreateTemp创建临时文件
- [ ] 异常情况下也要确保文件被删除

```go
// ✅ 正确示例
func processTemp() error {
    tmpFile, err := os.CreateTemp("", "prefix-")
    if err != nil {
        return err
    }
    defer os.Remove(tmpFile.Name())  // 确保删除
    defer tmpFile.Close()

    // 处理文件...
    return nil
}
```

### 2.6 文件上传检查 【强制】
- [ ] 文件保存时对文件名随机化处理
- [ ] 使用白名单验证文件后缀
- [ ] 验证文件MIME类型
- [ ] 限制文件大小
- [ ] 将文件存储在非Web目录或配置禁止执行

```go
// ✅ 正确示例
func saveUploadedFile(file multipart.File, header *multipart.FileHeader) error {
    // 验证文件类型
    ext := filepath.Ext(header.Filename)
    if !isAllowedExtension(ext) {
        return errors.New("file type not allowed")
    }

    // 生成随机文件名
    randomName := generateRandomFilename() + ext
    savePath := filepath.Join(uploadDir, randomName)

    // 保存文件...
    return nil
}
```

---

## 03. 系统接口

### 3.1 禁止调用OS命令解析器或运行程序防止命令注入 【强制】
- [ ] 避免直接拼接用户输入到系统命令
- [ ] 使用exec.Command而非shell命令
- [ ] 不使用shell解析器（避免使用sh -c）
- [ ] 对命令参数进行严格的白名单校验
- [ ] 避免使用os/exec执行动态命令

```go
// ❌ 错误示例
func execCommand(userInput string) {
    cmd := exec.Command("sh", "-c", "ls "+userInput)  // 命令注入风险
    cmd.Run()
}

// ✅ 正确示例
func execCommand(dir string) error {
    // 严格校验输入
    if !isValidPath(dir) {
        return errors.New("invalid path")
    }

    // 直接使用命令，不通过shell
    cmd := exec.Command("ls", dir)
    return cmd.Run()
}
```

---

## 04. 通讯安全

### 4.1 网络通信采用TLS方式 【可选】
- [ ] 网络通信使用TLS加密传输
- [ ] gRPC/Websocket使用TLS 1.3
- [ ] 避免使用明文传输协议（HTTP、FTP等）
- [ ] 配置强加密套件

### 4.2 TLS启用证书验证 【可选】
- [ ] 生产环境必须启用证书验证
- [ ] 证书应当有效、未过期
- [ ] 证书配置正确的域名
- [ ] 不使用InsecureSkipVerify

```go
// ❌ 错误示例
tr := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},  // 不安全
}

// ✅ 正确示例
tr := &http.Transport{
    TLSClientConfig: &tls.Config{
        MinVersion: tls.VersionTLS13,
        // 启用证书验证（默认）
    },
}
```

### 4.3 Websocket需要校验Origin 【可选】
- [ ] Websocket连接时校验Origin头
- [ ] 使用白名单验证Origin
- [ ] 防止CSRF攻击

```go
// ✅ 正确示例
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        return isAllowedOrigin(origin)
    },
}
```

---

## 05. 敏感数据保护

### 5.1 敏感信息访问 【强制】
- [ ] 禁止将敏感信息硬编码在程序中
- [ ] 敏感信息从配置文件或环境变量读取
- [ ] 密钥、密码使用密钥管理服务
- [ ] 配置文件权限设置为仅所有者可读（0600）

```go
// ❌ 错误示例
const apiKey = "sk-1234567890abcdef"  // 硬编码

// ✅ 正确示例
func getAPIKey() (string, error) {
    return os.Getenv("API_KEY"), nil  // 从环境变量读取
}
```

### 5.2 敏感信息输出 【强制】
- [ ] 只输出必要的最小数据集
- [ ] 避免在响应中包含敏感字段
- [ ] 日志中不记录敏感信息（密码、密钥、身份证等）
- [ ] 错误信息不暴露系统细节

```go
// ❌ 错误示例
log.Printf("User login: username=%s, password=%s", user, pass)

// ✅ 正确示例
log.Printf("User login: username=%s", user)  // 不记录密码
```

### 5.3 敏感数据存储 【强制】
- [ ] 敏感数据使用SHA256、RSA等算法加密存储
- [ ] 密码使用bcrypt、scrypt、Argon2等慢哈希算法
- [ ] 不使用MD5、SHA1等弱哈希算法存储密码
- [ ] 使用加盐哈希（salt + hash）

```go
// ✅ 正确示例
import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
```

### 5.4 异常处理和日志记录 【强制】
- [ ] 合理使用panic、recover、defer处理异常
- [ ] 避免错误信息输出到前端
- [ ] 前端只返回友好的错误提示
- [ ] 详细错误信息记录到日志

```go
// ✅ 正确示例
func handleRequest(w http.ResponseWriter, r *http.Request) {
    err := processRequest(r)
    if err != nil {
        log.Printf("Error processing request: %v", err)  // 详细信息记录到日志
        http.Error(w, "Internal Server Error", 500)      // 前端返回通用错误
        return
    }
}
```

---

## 06. 加密解密

### 6.1 不得硬编码密码/密钥 【强制】
- [ ] 禁止在代码中硬编码密码或密钥
- [ ] 密钥从配置文件、环境变量或密钥管理服务获取
- [ ] 可通过变换算法保护密钥
- [ ] 定期轮换密钥

```go
// ❌ 错误示例
const encryptKey = "1234567890abcdef"  // 硬编码

// ✅ 正确示例
func getEncryptKey() ([]byte, error) {
    key := os.Getenv("ENCRYPT_KEY")
    if key == "" {
        return nil, errors.New("encrypt key not set")
    }
    return []byte(key), nil
}
```

### 6.2 密钥存储安全 【强制】
- [ ] 使用对称加密时保护好密钥
- [ ] 涉及敏感数据时使用非对称算法协商密钥
- [ ] 密钥文件权限设置为0600
- [ ] 使用硬件安全模块（HSM）或密钥管理服务（KMS）

### 6.3 不使用弱密码算法 【可选】
- [ ] 不使用DES、3DES、RC4等弱加密算法
- [ ] 不使用MD5、SHA1等弱哈希算法
- [ ] 使用AES-256、ChaCha20等强加密算法
- [ ] 使用SHA256、SHA512等强哈希算法

```go
// ❌ 错误示例 - 使用DES
import "crypto/des"

// ✅ 正确示例 - 使用AES
import "crypto/aes"
import "crypto/cipher"

func encrypt(plaintext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    // 使用GCM模式...
}
```

### 6.4 使用crypto/rand包生成安全随机数 【强制】
- [ ] 使用crypto/rand生成密钥、token等安全随机数
- [ ] 不使用math/rand生成安全相关的随机数
- [ ] math/rand的随机数是可预测的

```go
// ❌ 错误示例
import "math/rand"
token := rand.Int63()  // 不安全

// ✅ 正确示例
import "crypto/rand"

func generateToken() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}
```

---

## 07. 正则表达式

### 7.1 使用regexp进行正则表达式匹配 【可选】
- [ ] 使用regexp包进行正则匹配（保证线性时间性能）
- [ ] 防止正则表达式DoS攻击（ReDoS）
- [ ] 对正则表达式进行性能测试
- [ ] 设置正则匹配超时时间

```go
// ✅ 正确示例
import "regexp"

func validateInput(input string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9_-]{1,50}$`)
    return re.MatchString(input)
}
```

---

## 08. 输入输出校验

### 8.1 按类型进行数据校验 【强制】
- [ ] 所有外部输入使用validator进行白名单校验
- [ ] 校验内容：数据长度、数据范围、数据类型、数据格式
- [ ] 校验不通过时拒绝请求
- [ ] 使用validator库进行结构化校验

```go
// ✅ 正确示例
import "github.com/go-playground/validator/v10"

type User struct {
    Username string `validate:"required,alphanum,min=3,max=20"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"required,gte=0,lte=130"`
}

func validateUser(user *User) error {
    validate := validator.New()
    return validate.Struct(user)
}
```

---

## 09. SQL操作

### 9.1 SQL语句默认使用预编译并绑定变量 【强制】
- [ ] 使用参数化查询（Prepared Statement）
- [ ] 禁止拼接SQL语句
- [ ] 使用ORM框架时确保使用参数化查询
- [ ] 防止SQL注入攻击

```go
// ❌ 错误示例
query := "SELECT * FROM users WHERE username = '" + username + "'"  // SQL注入风险
db.Query(query)

// ✅ 正确示例
query := "SELECT * FROM users WHERE username = ?"
db.Query(query, username)  // 使用参数化查询
```

---

## 10. 网络请求

### 10.1 资源请求过滤验证 【强制】
- [ ] 使用http.Get/Post/Do时，对URL进行严格校验
- [ ] 外部可控URL必须使用白名单验证
- [ ] 防止SSRF（服务端请求伪造）攻击
- [ ] 禁止访问内网地址（127.0.0.1、192.168.x.x等）
- [ ] 设置请求超时时间

```go
// ✅ 正确示例
func fetchURL(urlStr string) error {
    // 解析URL
    u, err := url.Parse(urlStr)
    if err != nil {
        return err
    }

    // 验证协议
    if u.Scheme != "http" && u.Scheme != "https" {
        return errors.New("invalid protocol")
    }

    // 验证主机不是内网地址
    if isInternalIP(u.Host) {
        return errors.New("internal IP not allowed")
    }

    // 设置超时
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(urlStr)
    // ...
}
```

---

## 11. 服务器端渲染

### 11.1 模板渲染过滤验证 【强制】
- [ ] 使用text/template或html/template渲染模板
- [ ] 禁止将外部输入直接引入模板
- [ ] 仅允许引入白名单内字符
- [ ] html/template会自动转义，text/template需手动转义
- [ ] 防止模板注入攻击

```go
// ✅ 正确示例
import "html/template"

func renderPage(w http.ResponseWriter, data interface{}) {
    tmpl := template.Must(template.ParseFiles("page.html"))
    tmpl.Execute(w, data)  // html/template会自动转义
}
```

---

## 12. Web跨域

### 12.1 跨域资源共享CORS限制请求来源 【可选】
- [ ] 严格设置Access-Control-Allow-Origin
- [ ] 使用同源策略保护
- [ ] 避免设置为*（允许所有来源）
- [ ] 使用白名单验证Origin

```go
// ✅ 正确示例
func handleCORS(w http.ResponseWriter, r *http.Request) {
    origin := r.Header.Get("Origin")
    if isAllowedOrigin(origin) {
        w.Header().Set("Access-Control-Allow-Origin", origin)
    }
    // 不要使用: w.Header().Set("Access-Control-Allow-Origin", "*")
}
```

---

## 13. 响应输出

### 13.1 设置正确的HTTP响应包类型 【强制】
- [ ] Content-Type与实际响应内容一致
- [ ] JSON响应使用application/json
- [ ] XML响应使用text/xml
- [ ] HTML响应使用text/html

```go
// ✅ 正确示例
func sendJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
```

### 13.2 添加安全响应头 【可选】
- [ ] 添加X-Content-Type-Options: nosniff
- [ ] 添加X-Frame-Options: DENY或SAMEORIGIN
- [ ] 添加X-XSS-Protection: 1; mode=block
- [ ] 添加Content-Security-Policy
- [ ] 添加Strict-Transport-Security（HTTPS）

```go
// ✅ 正确示例
func setSecurityHeaders(w http.ResponseWriter) {
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "DENY")
    w.Header().Set("X-XSS-Protection", "1; mode=block")
    w.Header().Set("Content-Security-Policy", "default-src 'self'")
}
```

### 13.3 外部输入拼接到HTTP响应头中需进行过滤 【强制】
- [ ] 避免外部可控参数拼接到HTTP响应头
- [ ] 必须拼接时过滤\r、\n等换行符
- [ ] 拒绝携带换行符的外部输入
- [ ] 防止HTTP响应拆分攻击

```go
// ✅ 正确示例
func setCustomHeader(w http.ResponseWriter, value string) error {
    // 检查是否包含换行符
    if strings.ContainsAny(value, "\r\n") {
        return errors.New("invalid header value")
    }
    w.Header().Set("X-Custom-Header", value)
    return nil
}
```

### 13.4 外部输入拼接到response页面前进行编码处理 【强制】
- [ ] 使用html/template自动编码
- [ ] 使用html.EscapeString编码<, >, &, ', "
- [ ] 防止XSS攻击

```go
// ✅ 正确示例
import "html/template"

func renderUserInput(w http.ResponseWriter, userInput string) {
    tmpl := template.Must(template.New("page").Parse(`<p>{{.}}</p>`))
    tmpl.Execute(w, userInput)  // 自动转义
}
```

---

## 14. Session安全

### 14.1 安全维护session信息 【强制】
- [ ] 用户登录时重新生成session ID
- [ ] 退出登录后清理session
- [ ] session ID长度足够（至少128位）
- [ ] session ID使用安全随机数生成
- [ ] 设置合理的session过期时间
- [ ] session cookie设置HttpOnly和Secure标志

```go
// ✅ 正确示例
func login(w http.ResponseWriter, r *http.Request) {
    // 登录成功后重新生成session
    session, _ := store.Get(r, "session-name")
    session.Options = &sessions.Options{
        MaxAge:   3600,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    }
    session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    session.Options.MaxAge = -1  // 清理session
    session.Save(r, w)
}
```

### 14.2 CSRF防护 【强制】
- [ ] 涉及敏感操作的接口校验Referer
- [ ] 添加csrf_token验证
- [ ] 使用框架默认的CSRF防御机制
- [ ] token使用安全随机数生成

```go
// ✅ 正确示例
import "github.com/gorilla/csrf"

func main() {
    CSRF := csrf.Protect(
        []byte("32-byte-long-auth-key"),
        csrf.Secure(true),
    )

    http.Handle("/", CSRF(handler))
}
```

---

## 15. 访问控制

### 15.1 默认鉴权 【强制】
- [ ] 系统默认进行身份认证
- [ ] 使用白名单方式放开不需要认证的接口
- [ ] 所有接口默认需要鉴权，除非明确可公开访问
- [ ] 实现统一的认证中间件

```go
// ✅ 正确示例
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 白名单检查
        if isPublicPath(r.URL.Path) {
            next.ServeHTTP(w, r)
            return
        }

        // 认证检查
        if !isAuthenticated(r) {
            http.Error(w, "Unauthorized", 401)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

---

## 16. 并发保护

### 16.1 禁止在闭包中直接调用循环变量 【强制】
- [ ] 循环中启动goroutine时，避免直接使用循环变量
- [ ] 使用参数传递或创建局部变量
- [ ] 防止数据竞争

```go
// ❌ 错误示例
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i)  // 所有goroutine可能打印相同的值
    }()
}

// ✅ 正确示例
for i := 0; i < 10; i++ {
    i := i  // 创建局部变量
    go func() {
        fmt.Println(i)
    }()
}

// 或使用参数传递
for i := 0; i < 10; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)
}
```

### 16.2 禁止并发写map 【强制】
- [ ] 并发访问map时使用锁保护
- [ ] 使用sync.Map处理并发map访问
- [ ] 避免程序崩溃

```go
// ❌ 错误示例
var m = make(map[string]int)
go func() { m["key"] = 1 }()  // 并发写入会崩溃
go func() { m["key"] = 2 }()

// ✅ 正确示例 - 使用mutex
var (
    m  = make(map[string]int)
    mu sync.Mutex
)

func setMap(key string, value int) {
    mu.Lock()
    defer mu.Unlock()
    m[key] = value
}

// 或使用sync.Map
var m sync.Map
m.Store("key", 1)
```

### 16.3 确保并发安全 【强制】
- [ ] 敏感操作进行并发安全限制
- [ ] 使用sync.Mutex或sync.RWMutex保护共享资源
- [ ] 使用原子操作（sync/atomic）
- [ ] 使用channel进行协程间通信
- [ ] 防止数据竞争

```go
// ✅ 正确示例 - 使用atomic
import "sync/atomic"

var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}

func getCounter() int64 {
    return atomic.LoadInt64(&counter)
}
```

---

## 17. 异常行为

### 17.1 禁止在异常中泄露敏感信息 【强制】
- [ ] 异常信息不包含敏感数据
- [ ] 异常信息不暴露系统实现细节
- [ ] 堆栈信息只记录到日志，不返回给用户
- [ ] 前端只返回通用错误信息

```go
// ❌ 错误示例
if err != nil {
    http.Error(w, err.Error(), 500)  // 可能泄露敏感信息
}

// ✅ 正确示例
if err != nil {
    log.Printf("Error: %v", err)  // 详细信息记录到日志
    http.Error(w, "Internal Server Error", 500)  // 通用错误给用户
}
```

### 17.2 确保方法异常时对象能恢复到之前的状态 【强制】
- [ ] 方法失败时保证对象状态一致性
- [ ] 使用defer恢复资源状态
- [ ] 实现事务性操作

```go
// ✅ 正确示例
func (s *Service) UpdateConfig(newConfig Config) error {
    oldConfig := s.config
    s.config = newConfig

    if err := s.validate(); err != nil {
        s.config = oldConfig  // 恢复到之前的状态
        return err
    }

    return nil
}
```

### 17.3 确保异常情况下能释放文件句柄、DB连接和内存资源等 【强制】
- [ ] 使用defer确保资源释放
- [ ] 异常情况下也要释放资源
- [ ] 使用context控制超时

```go
// ✅ 正确示例
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // 确保关闭

    // 处理文件...
    return nil
}

func queryDB(ctx context.Context) error {
    conn, err := db.Conn(ctx)
    if err != nil {
        return err
    }
    defer conn.Close()  // 确保关闭

    // 查询数据库...
    return nil
}
```

---

## 18. 序列化与反序列化

### 18.1 禁止序列化未加密的敏感数据 【强制】
- [ ] 序列化前对敏感数据加密
- [ ] 使用json:"-"标签避免序列化敏感字段
- [ ] 实现自定义MarshalJSON方法过滤敏感字段

```go
// ✅ 正确示例
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"`  // 不序列化密码
    Email    string `json:"email"`
}
```

### 18.2 确保将敏感对象发送出信任区域前必须先签名后加密 【强制】
- [ ] 发送敏感数据前先签名再加密
- [ ] 验证数据完整性（签名验证）
- [ ] 保证数据机密性（加密）
- [ ] 接收方先解密后验签

```go
// ✅ 正确示例
func sendSensitiveData(data []byte) ([]byte, error) {
    // 1. 签名
    signature, err := sign(data)
    if err != nil {
        return nil, err
    }

    // 2. 组合数据和签名
    signedData := append(data, signature...)

    // 3. 加密
    encryptedData, err := encrypt(signedData)
    if err != nil {
        return nil, err
    }

    return encryptedData, nil
}
```

---

## 附录：安全编码工具

### 推荐使用的安全库

#### 密码学
```go
// 加密
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"crypto/sha256"
"golang.org/x/crypto/bcrypt"  // 密码哈希

// TLS
"crypto/tls"
```

#### 数据验证
```go
"github.com/go-playground/validator/v10"  // 数据校验
"github.com/microcosm-cc/bluemonday"      // HTML清理
```

#### Web安全
```go
"github.com/gorilla/csrf"        // CSRF防护
"github.com/gorilla/securecookie" // 安全Cookie
"golang.org/x/crypto/acme/autocert" // 自动TLS证书
```

#### SQL安全
```go
// 使用ORM框架
"gorm.io/gorm"
"github.com/jmoiron/sqlx"
```

---

## 安全编码检查清单

### 代码提交前检查
- [ ] 运行golangci-lint进行静态代码检查
- [ ] 运行go test -race检查数据竞争
- [ ] 运行gosec进行安全扫描
- [ ] 检查是否有硬编码的敏感信息
- [ ] 检查所有外部输入是否进行了校验
- [ ] 检查错误处理是否完整

### 代码审查重点
- [ ] 输入验证和输出编码
- [ ] 认证和授权机制
- [ ] 敏感数据处理
- [ ] 错误处理和日志记录
- [ ] 并发安全
- [ ] 资源管理

### 安全测试
- [ ] 进行渗透测试
- [ ] 进行SQL注入测试
- [ ] 进行XSS测试
- [ ] 进行CSRF测试
- [ ] 进行SSRF测试
- [ ] 进行竞态条件测试

---

## 推荐工具

### 静态分析工具
```bash
# 安全扫描
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...

# 数据竞争检测
go test -race ./...

# 代码质量检查
golangci-lint run
```

### CI/CD集成
```yaml
# .github/workflows/security.yml
name: Security Scan
on: [push, pull_request]
jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: ./...
```

---

## 参考资源

- [OWASP Go Secure Coding Practices](https://owasp.org/www-project-go-secure-coding-practices-guide/)
- [Go Security Best Practices](https://golang.org/doc/security/)
- [CWE Top 25](https://cwe.mitre.org/top25/)

---

*最后更新: 2025-11-24*
*版本: v1.0*
