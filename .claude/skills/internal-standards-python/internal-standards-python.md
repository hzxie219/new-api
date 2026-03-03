| # | Checklist 项 |
|:---|:---|
| | 01. style - 风格规范 |
| 1.1. | label - 标识符命名：<br>1) 模块与包命名尽量短小，全部使用小写，模块命名允许使用下划线，使用名词。【强制】<br>2）类命名采用名词，首字母大写，多个词组合时，每个词的首字母大写。【强制】<br>3) 异常命名使用CapWords+Error后缀的方式。【强制】<br>4）函数与方法命名采用动宾结构，全部使用小写，各个词中间采用下划线隔开,不使用少于三个字符的函数名字。【强制】<br>5）变量(包括全局变量、类变量、实例变量、句柄变量)、函数参数、方法参数使用全小写加下划线，不使用单个字符的变量名。【强制】<br>   例外场景：<br>   A) 计数器和迭代器 (例如, i, j, k, v 等等)。<br>   B) 在try/except语句中代表异常的e。<br>   C) 在with语句中代表文件句柄的f。<br>   D) 私有的、没有约束的类型变量 (type variable）, 例如_T = TypeVar("_T"), _P = ParamSpec("_P")。<br>6) 常量命名使用全部大写的方式，可以使用下划线。【强制】<br>7) 文件名都应该以.py为文件后缀且不能包含连字符 (-)。【强制】<br>8) 类型别名应该采用CapWords方式。【强制】<br>9) 类型变量应该采用CapWords方式，并且应使用短名称：T，AnyStr。如果是声明协变量或逆变量行为的变量，应添加_co或_contra后缀。【强制】<br>9）变量名字不能隐藏内置名字。【强制】<br>10) 在同一代码文件中，命名方式需统一。【强制】<br>11) 有文档辅助的接口，比如sdk，命名应当短小精练，不能带类型，在业务流程代码，缺少文档的命名，比如局部变量，文件全局变量，应该有类型标识。【建议】<br>12) 类的属性若与关键字名字冲突，后缀一下划线，尽量不要使用缩略等其他方式。【建议】<br>13) 为避免与子类属性命名冲突，在类的一些属性前，前缀两条下划线。比如：类Foo中声明__a,访问时，只能通过Foo._Foo__a，避免歧义。如果子类也叫Foo，那就无能为力了。【建议】<br>14) 实例方法的第一个参数必须为self,类方法的第一个参数必须为cls。【建议】 |
| 1.2. | indent - 缩进：<br>缩进采用4个空格的缩进（编辑器都可以完成此功能），不使用Tap，更不能混合使用Tap和空格。【强制】 |
| 1.3. | blank - 空行：<br>类和top-level函数定义之间空两行；类中的方法定义之间空一行；函数内逻辑无关段落之间空一行；if、while、for等控制语句之后空一行；其他地方尽量不要再空行。【强制】 |
| 1.4. | wrap - 折行：<br>1) 续行应采用圆括号, 中括号和花括号的隐式续行,不要用反斜杠表示显式续行。【强制】<br>   如有需要,可以在表达式外围添加一对括号。<br>   例外：3.9之前的版本的with语句。<br>2) 括号内换行采用垂直隐式连接&悬挂式缩进。使用悬挂式缩进应该注意第一行不应该有参数，连续行要使用进一步的缩进来区分。【建议】<br>   A) 右括号 (圆括号, 方括号或花括号) 可以置于表达式结尾或者另起一行。另起一行时右括号应该和左括号所在的那一行缩进相同。<br>   B) 换行时需要防止换行的语句和后续的代码缩进相同。<br>3) 长表达式要在低优先级操作符处拆分续行，操作符放在新行之首（以便突出操作符）。【建议】 |
| 1.5. | parenthesis - 括号：<br>使用括号时宁缺毋滥。不要在返回语句或条件语句中使用括号, 除非用于隐式续行或表示元组。【强制】 |
| 1.6. | comment - 注释：<br>1) 采用"#"在模块前面进行版本许可声明。【建议】<br>2) 模块头：采用模块文档字符串(docstring)，需注释说明该文件的用途和注意事项。【建议】<br>3) 函数头：函数文档字符串应注释(docstring)，需说明其功能及各参数、返回值的含义。【建议】<br>4) 类头：类采用文档注释(docstring)，需说明类的功能；【建议】<br>5) 文件头：在文件头部注释该文件的功能和作者信息，必须注释该文件的所有作者和修改者，并说明各作者和修改者的修改内容，中文注释格式为"作者：xxx","修改者: xxx"，英文注释格式为"author: xxx","mender: xxx"；【建议】<br>6) 行注释需要在行上面做注释，不推荐行尾注释；【建议】<br>   注释的每一行都以#号和一个空格开始，缩进与该代码相同的水平。<br>7) 在下述控制结构处应按要求进行注释。语句块少于5行允许例外。【建议】<br>   A) if语句的各个分支，注释说明条件和具体功能；<br>   B) for/while的头部，注释说明循环条件和具体功能；<br>8) 优先使用中文注释，除非英语水平很高；【建议】 |
| 1.7. | space - 空格：<br>总体原则，避免不必要的空格。<br>1) 各种右括号前不要加空格。【强制】<br>2) 逗号、冒号、分号前不要加空格。【强制】<br>3) 函数的左括号前不要加空格。如Func(1)。【强制】<br>4) 序列的左括号前不要加空格。如list[2]。【强制】<br>5) 操作符左右各加一个空格，不要为了对齐增加空格。【强制】<br>6) 函数默认参数使用的赋值符左右省略空格。【强制】<br>7) 函数注解中，冒号与参数之间无空格，与类型之间加一个空格；箭头前后都加一个空格。【强制】<br>8) 函数默认参数和参数注解一起使用时，赋值符前后都加一个空格。【强制】 |
| 1.8. | trailing_comma - 尾随逗号：<br>1) 仅当],),}和最后一个元素不在同一行时, 推荐在序列尾部添加逗号。【建议】 |
| 1.9. | coding_format - 编码格式：<br>1) 代码文件编码格式必须为"UTF-8无BOM编码格式"。【强制】<br>2) 新增文件必须增加注释 # coding=utf-8。【强制】<br>3) 每个模块都应该以#!/usr/bin/env python<version>开头。【强制】<br>4) 换行符必须采用unix格式的换行符。【强制】<br>5) 禁止使用sys.setdefaultencoding方法来设置编码。【强制】 |
| | 02. exception - 错误处理 |
| 2.1. | raise - 异常抛出：<br>1) 对于处理从外部读取的数据时，比如文件、数据库、标准输入、命令行参数、管道、socket、及其它任意RPC机制，对可能出现异常的地方增加try/except处理，防止程序的异常退出。【强制】<br>2) 抛出异常时，使用raise ValueError('message')，不要使用两个参数形式（raise MyException, "Error message"）或者过时的字符串异常（raise "Error message"）【强制】<br>3) except和finally子句中，不得抛出新的未知异常，防止覆盖当前的异常。python3里，except里需要抛新异常，建议采用异常链方式：raise X from Y。【建议】<br>   except作用域中不能进行复杂的代码逻辑，最好只调用LOG.exception打日志和异常类型转换，如果需要进行如资源回滚删除等操作，需要用try捕获异常，避免覆盖当前的异常。 |
| 2.2. | catch - 异常捕获：<br>1) 不要使用except:语句捕获所有异常，也不要捕获Exception或者StandardError。【强制】<br>   例外：当前处理逻辑已经处于最外层或者捕获到异常之后重新触发该异常。<br>2) 异常必须同级的try catch捕获，禁止在同级的捕获语句捕获该异常（实际是捕获不到的）【强制】<br>3) 捕获多个异常时，使用except (Exception1, Exception2):。【强制】 |
| 2.3. | exception - 异常处理：<br>1) 使用基于对象的异常，不要使用已废弃的字符串异常。【强制】<br>2) 模块或者包应当定义自己的异常基类，这个类应当继承自内置的Exception类。【强制】<br>3) 将无论异常与否都应执行的代码放在finally里。【强制】<br>   这种写法常用于清理资源, 例如关闭文件。<br>4) 严禁在finally中使用流控制语句return/break/continue。【强制】<br>5) 严禁在exception里通过pass忽略异常，如果需要，必须进行审核并打印日志。【强制】<br>6) 尽量减少try/except块中的代码量，try块的体积越大，期望之外的异常就越容易被触发，此时try/except会隐藏真正的错误。【建议】 |
| 2.4. | log - 异常日志规范:<br>1) LOG.exception只能用于发现异常处的except作用域内。【强制】<br>2) 所有异常的地方都需要用LOG.exception打印堆栈，不能使用LOG.error，除非能确定异常抛出地方的除外，如一些约定好的异常类型。【建议】<br>3) LOG.exception打印时不需要将异常对象e放在LOG中，并且需要把出现异常的条件和相关变量打印出来。【建议】 |
| 2.5. | assert - 断言要求：<br>1) 优先使用合适的内置异常类。assert只能用来保证内部正确性，不应该用assert来纠正参数或表示意外情况。【强制】<br>2) 断言中禁止调用有副作用的函数。【强制】<br>   有副作用的函数，指函数内会更改变量值、改变系统环境、进行IO。 |
| | 03. func - 函数/方法: |
| 3.1. | default_parameter - 函数参数默认值：<br>1）禁止通过默认参数来实现每次调用都获取不同值，如now=time.time()。【强制】<br>2）不要使用可变对象作为函数默认值，如列表或字典。【强制】<br>3) 函数默认参数赋值不能调用业务逻辑函数。【强制】<br>4) 新添加的参数须考虑向后兼容问题，比如新增参数带有效默认值。【强制】 |
| 3.2. | var_parameter - 可变参数：<br>尽量不使用foo(*args, **kwargs)这样的写法。【建议】 |
| 3.3. | cmd_parameter - 命令行参数：<br>命令行参数的值不允许使用json格式。【强制】 |
| 3.4. | parameter - 参数：<br>禁止使用locals()，globals()传递参数。【强制】 |
| 3.5. | return - 返回值：<br>1) 函数返回值必须小于等于3个。3个以上时必须通过class/namedtuple/dict等具名形式进行包装。【强制】<br>2) 返回值如需要判断函数运行的正确性，不通过返回值来检查，而是通过抛异常的方式。【建议】<br>3) 返回值应保证一致性，不一致的返回值会造成上层调用需要做特判。【建议】<br>   A) 如果有返回值，那么return需携带明确的返回值，如果没有则明确return None，并且在函数的结尾应该有一个明确的返回语句。<br>   B) 如果有返回值，返回值类型需要保持一致。 |
| 3.6. | properties - 特性：<br>1) 使用@property装饰器来创建特性，不能自行实现特性装饰器。【强制】<br>2) 在继承场景中，不要在子类中覆写或扩展父类的特性实现。【强制】<br>3) 特性在读取或设置属性时，不能仅仅只是获取和设置一个内部属性，应涉及简单的计算逻辑。【建议】 |
| 3.7. | type_annotation - 代码类型注解：<br>对python3代码进行代码类型注解,并使用诸如pytype之类的类型检查工具来检查代码，特别是泛型类型。【建议】 |
| 3.8. | internal_fun - 嵌套/局部/内部类或函数：<br>不推荐使用嵌套/局部/内部类或函数，在不增加理解成本的情况下，可以适当例外。【建议】 |
| 3.9. | lexical_scope - 词法作用域：<br>在必须使用嵌套函数的场景下<br>1) 禁止在闭包中直接进行变量绑定。【强制】<br>2) 嵌套函数只允许只读访问外层作用域中定义的变量，并且外部作用域里的变量必须进行定义。【强制】 |
| 3.10. | lambda - lambda函数：<br>1) 对于常见的操作符，例如乘法操作符，使用operator模块中的函数以代替lambda函数。【强制】<br>   例如, 推荐使用operator.mul, 而不是lambda x, y: x * y。<br>2) 适用于单行函数. 如果代码超过60-80个字符, 最好还是定义成常规(嵌套)函数。【建议】 |
| 3.11. | decorator - 函数与方法装饰器：<br>1) 实现装饰器时，返回的函数需用six.wraps(func)或functools.wraps(func)装饰，func为被装饰函数。【强制】<br>2) 避免装饰器自身对外界的依赖(即不要依赖于文件, socket, 数据库连接等)，也不能在装饰器中调用业务逻辑。【强制】<br>3) 如果好处很显然, 就明智而谨慎的使用装饰器. 装饰器应该遵守和函数一样的导入和命名规则. 装饰器的python文档应该清晰的说明该函数是一个装饰器. 请为装饰器编写单元测试.【建议】 |
| | 04. language - 语言特性 |
| 4.1. | import - 导入：<br>1) 导入始终放在文件顶部，紧随任何模块注释和文档字符串之后，以及模块全局变量和常量之前。【强制】<br>2) 在引用外部包的时候，禁止使用import *、from A import *； 【强制】<br>3）每个导入应该独占一行，禁止使用import os, sys; 【强制】<br>4）不要导入不使用的模块 【强制】<br>5）使用绝对导入，禁止使用相对导入。【强制】<br>   即使模块在同一个包中, 也要使用完整包名. 这能帮助你避免无意间导入一个包两次 <br>6）在引用外部包的时候，按照以下顺序排列，并在每组导入之间，插入空白行: 【建议】<br>    __future__<br>    标准库<br>    第三方库<br>    自己的库 |
| 4.2. | main - 主程序：<br>1) 代码必须在执行主程序前检查if __name__ == '__main__'。【强制】<br>2) 除了__main__的文件，其他python文件应禁止业务函数直接执行。【强制】 |
| 4.3. | dunders - 特殊变量:<br>对于两个_开头和结尾的变量，如__all__，__author__，__version__等，应该放在模块文档之后， 其他模块导入之前（__future__ 除外）。【建议】 |
| 4.4. | init_py - __init__.py：<br>__init__.py应该是空文件，禁止在__init__.py实现业务逻辑，可以适当允许使用导入控制。【强制】 |
| 4.5. | statement - 语句行：<br>1) 链式赋值制里，只允许应用于简单的变量赋值，禁止一个变量在链式赋值里多次使用。【强制】<br>2) 不要将多句语句写在同一行，尽管使用';'允许。【强制】<br>3) if/for/while语句中，即使执行语句只有一句，也必须另起一行。【强制】 |
| 4.6. | if - if语句：<br>1) 过长的if else语句，不能写成条件表达式，需要使用传统if else多条语句格式。【强制】<br>2) if的判断语句，建议不超过一行，可通过声明变量或提取函数保存判断条件来解决。【建议】<br>3) 尽量使用内建的all, any处理2个以上的条件判断，如if any((x, y, z)) 要优于if x or y or z。【建议】<br>4) if和else代码块的行数不能超过10行，超过需拆成子函数调用。【建议】<br>5) if优先判断处理逻辑较短的分支，处理逻辑复杂的分支最后的else或else if才处理。【建议】<br>6) if优先判断处理逻辑较短的分支，如果下面只有else一个分支，可以提前退出，后面无需再用else。【建议】<br>7) 流程处理就遵循短路原则，能提前返回或抛异常退出的代码写在前面。【建议】 |
| 4.7. | comprehensions - 列表推导：<br>简单的场景可以使用列表推导，其中mapping、loop、filter部分单独成行，且最多只能写一行，禁止多层loop或filter。【强制】 |
| 4.8. | iterator - 默认迭代器和操作符：<br>1) 对容器或文件的只读遍历，应该使用内置的迭代方法，不要使用返回list的方式遍历。【建议】<br>2) 对容器类型，使用in或not in判断元素是否存在。而不是has_key。【建议】 |
| 4.9. | generator - 生成器：<br>1) 当返回较长列表数据时建议使用yield和generator函数。【建议】<br>   解释:生成器函数可以避免返回大列表，占用过多内存、影响性能。同时还可以保持代码的优雅易读。<br>2) python2里，应使用xrange替换range生成数字序列。【建议】 |
| 4.10. | cond_expression - 条件表达式：<br>1) 条件表达式仅用于一行之内，禁止嵌套使用。【强制】 |
| 4.11. | bool - True/False的求值：<br>1) 尽可能使用隐式的false，0, None, [], {}, "" 作为布尔值使用时相当于 False。【强制】<br>   注意，在使用隐式false时，需要考虑是区分false和None。比如在处理整数时，可能需要区别0和None。<br>2) 不要使用==比较一个布尔值是否等于False。应该用if not x:或者if x is True:代替。如果需要区分False和None, 应该使用if not x and x is not None:这样的判断语句。【强制】<br>3) 判断序列(字符串, 列表, 元组)是否为空序列，因使用if not seq:（或if seq:）代替if len(seq):（或if not len(seq):）。【强制】<br>4) Numpy数组转换为布尔值时会抛出异常. 因此使用.size属性检查np.array是否为空。【强制】 |
| 4.12. | signletio_comparison - 单例对象比较：<br>使用is或者is not操作符来比较单例对象（比如None，True，False），而不是==或者!=运算符。【强制】<br>注意：<br>A) 不要使用not ... is操作符。<br>B) 不要使用if x:来检测None。 |
| 4.13. | obsolete - 过时的语言特性：<br>1) 使用字符串方法取代字符串模块。【强制】<br>2) 使用函数调用语法取代apply()。【强制】 |
| 4.14. | complex_feature - 威力过大的特性：<br>禁止使用威力过大的特性。【强制】<br>包括：<br>A）元类（metaclasses）<br>B）字节码访问<br>C）任意编译（on-the-fly compilation）<br>D）动态继承<br>E）对象父类重定义(object reparenting)<br>F）导入黑客（import hacks）<br>G）反射<br>H）系统内修改（modification of system internals） |
| 4.15. | string - 字符串：<br>1）多行字符串使用三重双引号，而不是三重单引号。【强制】<br>2) 多于2个字符串的拼接，禁止使用"+"，而应该用join。【强制】<br>3) 使用''.startswith() 和 ''.endswith()而非字符切片去检测前缀或后缀。【强制】<br>4) 应使用f-string、str.format()方式来格式化字符串，避免使用"%s " % (xx)，简单的格式化字符串可以适当例外。【强制】<br>5) 对于第一个参数是格式字符串 (包含 % 占位符) 的日志函数: 一定要用字符串字面量 (而非 f-string!) 作为第一个参数, 并用占位符的参数作为其他参数。【强制】<br>6) 禁止使用isalpha()判断字符串是否都是字母，应该使用正在表达式。【强制】<br>7）不要在循环中用+和+=操作符来累加字符串。由于字符串是不可变的，这样会创建不必要的临时对象，并且导致二次方的运行时间。【建议】 |
| 4.16. | dict - 字典：<br>1) 在确定存在key的情况下，推荐使用dictB['keyA']而不使用dictB.get('keyA')。【强制】<br>2) 所有使用get的地方都需要显式声明默认值，如dictA.get('a', None)，禁止使用dictA.get('a')。【强制】<br>3) 将字典对象保存到文件或者数据库时，应使用json.loads()、json.dumps()将字典转换为JSON字符串再保存。【强制】<br>4) 单字段更新用索引赋值，如d['k']=v。多字段更新用update，如d.update(k=v)。【强制】<br>5) 遍历字典时，不能假定输出是有序的，如想使用有序字典，可以使用OrderedDict。【强制】<br>6) 遍历字典时，应该使用for x in xxx形式，而不是for x in xxx.keys()形式遍历字典。【强制】<br>7) 对字典赋默认值建议使用setdefault()函数，如dictA.setdefault('keyA', 'A')，不使用if判断然后再赋值。【建议】<br>   注意：如果默认值通过函数或表达式生成的，可以考虑使用if判断。 |
| 4.17. | tuple - 元组：<br>在创建只有一个元素的元组时，必须加逗号。如：(1,)。【强制】 |
| 4.18. | global - 全局变量：<br>禁止使用全局变量，采用类加静态变量的方式实现。【强制】<br>例外场景：<br>A) 用于保存默认配置选项的场景，可适当例外。<br>B) 用于模块级常量用途的场景，如PI = 3.14159，可适当例外。<br>C) 用于提升程序效率或者简化代码的场景，如用于缓存，或者保存函数的返回值等，可适当例外。 |
| 4.19. | init - 变量初始化：<br>1）如果一定要使用全局变量，则全局变量需要进行初始化。【强制】<br>2）类成员变量需要在构造函数中进行初始化。【强制】<br>3）局部变量，需要在使用之前进行初始化。【强制】<br> |
| 4.20. | list_remove - 方法使用：<br>1) 在遍历列表或数组时，不允许删除列表或数组里面的元素。【强制】 |
| 4.21. | performance - 性能：<br>1) 编码应当适用于Python的多个实现。比如不要依赖CPython的高效字符串语句 a+=b，而应当使用join，从而保证在不同实现上的线性开销。【强制】<br>2) 对于可变对象（如列表，字典等），使用x += y代替x = x + y。【强制】<br>3) 处理复杂的字符串模式匹配和替换时，应使用re库，不应直接使用字符串操作。【强制】<br>4) 单个进程内部数据传递尽量不要使用dict，应该写attrs之类的数据类。【建议】<br>5) 对于经常执行in运算符，并且数据较大，不应使用列表和元祖。【建议】 |
| 4.22. | sort - 排序操作：<br>当使用复杂比较实现排序操作时，最好实现全部的六个比较操作。【建议】 |
| 4.23. | object_comparison - 对象比较：<br>对象类型比较必须使用isinstance()而非直接比较。【强制】 |
| 4.24. | reusable_integers - 可重用整数：<br>整数比较禁止使用is运算符。【强制】 |
| 4.25. | magic - 魔数：<br>避免使用魔数，用常量代替。【强制】 |
| 4.26. | resource - 资源管理：<br>1) 文件、套接字以及其他可关闭资源 (比如数据库连接) ，在使用完之后应显式地关闭它们。【强制】<br>   推荐使用"with"语句管理文件和类似的资源；对于不支持 with 语句且类似文件的对象, 应该使用contextlib.closing()。否则因增加注释说明如何管理资源的生命周期<br>2) 当某个类预计会创建大量的对象实例时，并且不会动态添加属性，则应使用__slots__来减少内存占用。【强制】<br>   注意：继承场景下，其子类如果不定义自己的__slots__，将无法再添加新的实例属性。 |
| 4.27. | time - 时间：<br>1) 不能用time.timezone或者dateutil.tz.tzlocal()获取时区，可以用dateutil.tz.gettz()获取时区。【强制】<br>2) UTC时间不能用mktime()转时间戳，要用calendar.timegm()。【强制】 |
| 4.28. | protobuf - protobuf:<br>禁止在proto里使用python关键字。如果必须要使用，则python中应该通过setattr或getattr的方式来使用。【建议】 |
| 4.29. | hash - 哈希函数：<br>对于持久化的数据的一致性校验，禁止使用python3内置hash函数。【强制】 |
| 4.30. | class_attr - 类成员变量：<br>1) 对类成员变量禁止通过实例对象来操作。【强制】<br>   比如：<br>   A) 在成员函数内部通过self.class_variable = value来修改类成员变量。<br>   B) 通过实例对象直接对类成员变量赋值，如instance1.class_variable = value。<br>2) 对类成员变量进行定义时初始化如果是采用的函数返回的变量，则函数内禁止直接执行业务逻辑【强制】 |
| 4.31. | popen - 调用外部程序：<br>使用subprocess32调用外部命令需加上参数close_fds=True，并且需对stderr进行重定向，比如stderr=subprocess.DEVNULL。【强制】 |
| 4.32. | uuid - uuid：<br>uuid在比较时，需全部都转成string类型，除非可以保证都是uuid类型。【强制】 |
| 4.33. | deserialize - 反序列化：<br>1) 不使用yaml_load反序列化执行程序【强制】<br>   建议:使用yaml.safe_load来处理未受信任的YAML序列化<br>2) 不使用pickle模块序列化【强制】<br>   反序列容易造成代码注入，并且可能存在python不同版本序列化格式兼容性问题。 |
| | 05. concurrent - 并发 |
| 5.1. | thread_coroutine - 线程/协程：<br>1) 当一个父任务（如一个异步函数）启动一个或多个协程时，在父任务结束时必须进行协程回收。【强制】<br>2) 禁止在协程内部调用阻塞函数。【强制】<br>3) python多线程有其特定用途,而且不像C一样引发问题，但python多线程实现有缺陷，建议不使用，并且CPU密集型的业务禁止使用。【建议】<br>4) 协程库：python2推荐使用tornado和twisted，谨慎使用gevent，python3推荐使用asyncio和tornado【建议】 |
| 5.2. | process - 进程：<br>进程之间的通信和数据共享，如果涉及对大对象的频繁更新，禁止使用multiprocessing库。【强制】 |
| 5.3. | flush - 文件刷新：<br>并发读写文件时，需要及时flush文件或者关闭文件。【强制】 |
| | 06. i18n - 国际化 |
| 6.1. | punctuation - 标点符号:<br>1）翻译中英文使用英文标点， 中文使用中文标点。【强制】<br>2）分隔同类对象时，中文用顿号，英文用逗号 【强制】 |
| 6.2. | translate - 翻译:<br>1）所有需要翻译的字符串都需要用_() 【强制】<br>2）翻译的词条不应为格式化后的字符串，如_(u'未知错误！errno: %s' % errno) 【强制】<br>3）需要翻译的字串中，禁止出现"%s"的占位符，使用"%(name)s"带有命名意义的占位符。原因：对翻译不友好，缺少上下文【强制】 |
| 6.3. | string_encoding - 字符串编码:<br>1）如果一个地方返回的字符串编码可能是非ascii， 则一律返回utf-8。【强制】<br>2）禁止代码出现非ascii字符。【强制】<br>3) 如果对一个字符串编码， 选择ascii或utf-8， 原则上禁止其他编码， 除非特殊需要。【强制】 |
| | 07. security - 安全特性 |
| 7.1. | file_permissions - 文件权限：<br>设置较高的文件权限，可以执行更多的函数和操作 【强制】<br>建议:尽量不要给文件设置过高的权限 如：<br>os.chmod(key_file, 0o777) <br>os.chmod('/etc/passwd', 0o227） |
| 7.2. | hardcoded - 硬编码问题：<br>1) 不建议硬编码对接口的绑定，之后不易对代码进行修改，且可以对软件进行逆向获取 【强制】<br>建议:禁止对 '0.0.0.0' 地址进行绑定，这样有可能会无意开启一个未受到保护的通信服务<br>2) 不建议对密码进行硬编码，之后不易对代码进行修改，且可以对软件进行逆向获取 【强制】<br>建议：禁止将代码中的密码部分写死 例如禁止出现if password == "root": <br>3) 不建议对生成密码的函数默认参数进行硬编码 【强制】<br>建议:禁止出现doLogin(password="blerg") password参数应该使用可变参数来指定如doLogin(password=pass) |
| 7.3. | tmp_directory - 临时文件：<br>恶意用户可能在程序创建文件前预测文件名并劫持临时文件 【强制】<br>例如使用tempfile并且不要忘记删除临时文件，例如f = open('/tmp/abc', 'w')，tmp_dirs: ['/tmp', '/var/tmp', '/dev/shm'] 不要在这几个目录中创建临时文件，并且创建的文件要记得删除 |
| 7.4. | root_user - 用户权限：<br>以root权限进行执行，权限过高，可以执行任意命令 【强制】<br>建议:如以下函数<br>ceilometer_utils.execute('gcc --version', run_as_root=True) <br>cinder_utils.execute('gcc --version', run_as_root=True) <br>neutron_utils.execute('gcc --version', run_as_root=True) <br>nova_utils.execute('gcc --version', run_as_root=True) <br>nova_utils.trycmd('gcc --version', run_as_root=True) <br>将最后的参数改成False |
| 7.5. | command_injection - 命令注入：<br>1) eval可以任意将字符串执行,导致命令注入,使用ast.literal_eval代替 【强制】<br>2) input会从标准输入读取并执行代码,用raw_input代替 【强制】<br>   不安全的函数输入需要充分校验,不安全命令注入的函数列表：<br>   evec(),eval(),os.system(),os.popen(),execfile(),input(),compile()<br>3) 禁止使用linux命令通配符 【强制】<br>   建议:避免下面的写法 <br>   os.system("/bin/tar xvzf * ") <br>   os.popen2('/bin/chmod * ') <br>   使用*会导致无法意料到的错误 |
| 7.6. | weak_cryptographic - 安全加密算法：<br>1) random这个模块中的大部分随机数伪随机数，不能用于安全加密,使用 os.urandom()或者 random模块中的SystemRandom类来实现【强制】<br>   bad = random.random()<br>   bad = random.randrange()<br>   bad = random.randint()<br>   bad = random.choice()<br>   bad = random.uniform()<br>   bad = random.triangular()<br>   good = os.urandom()<br>   good = random.SystemRandom()<br>2) 低于1024位的DSA密钥大小被认为是不安全的，不使用 【强制】<br>   推荐的RSA和DSA算法的密钥长度大小大于2048 dsa.generate_private_key(key_size=2048,backend=backends.default_backend()) |
| 7.7. | xml - xml库：<br>在解析恶意构造的数据会产生一定的安全隐患,不使用系统自带的xml库，建议使用使用defusedxml库 【强制】 |
| 7.8. | ssl - ssl版本和认证：<br>1) request使用时未开启ssl认证,verify参数设为True 如：requests.get('https://gmail.com', verify=True)【强制】<br>2) 使用默认的SSL版本可能会存在漏洞【强制】<br>   //不安全的用法,使用了低版本的ssl协议 <br>   ssl.wrap_socket(ssl_version=ssl.PROTOCOL_SSLv2) <br>   //安全的用法,使用高版本的ssl协议 <br>   ssl.wrap_socket(ssl_version=ssl.PROTOCOL_TLSv1_2) |
| 7.9. | sql_injection - sql注入：<br>不使用sql字符串拼接的查询语句.【强制】<br>建议:要避免出现类似"SELECT %s FROM derp;" % var"SELECT thing FROM " + tab "SELECT " + val + " FROM " + tab + ... "SELECT {} FROM derp;".format(var)这样的语句 |
| 7.10. | frame - 框架相关：<br>flask:<br>1) 设置flask模板jinja2的autoescape为false，关闭自动转义，可能造成xss等相关漏洞【强制】<br>   建议:设置为true Environment(loader=templateLoader, load=templateLoader, autoescape=True)<br>2) 不能通过app.run(debug=True)来开启flask调试模式【强制】<br>Django:<br>3) 不能使用Django make_safe方法进行字符串转义 【强制】<br>   使用make_safe有可能会导致 xss等相关问题。所以不要使用这个函数，应手动对字符串进行转码。 |
| 7.11. | private_property - 私有属性：<br>1) 在变量名或函数名前加上"__"两个下划线来定义私有变量或函数。【建议】<br>   Python中默认的成员函数，成员变量都是公开的(public),而且python中没有类似public,private等关键词来修饰成员函数，成员变量。 |
| | 08. third_party_library - 第三方库检查 |
| 8.1. | check - 第三方库检查：<br>1) 如果对第三方库有改动,改动的代码需要遵守此编码规范。【强制】<br>2) 第三方库中未改动的代码需要进行代码扫描,风格和语言方面的问题可以不用修改,安全特性和bug类型的问题需要修改。【强制】 |
| | 09. tools - 工具检查 |
| 9.1. | check - 工具检查：<br>1) 建议尽量使用相关工具对代码检查，如代码扫描平台：flake8、pylint、pychecker。【建议】 |