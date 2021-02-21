## 1.UTXO结构封装
1.将所有与output相关的属性封装到一个结构
2.修改UTXO查找函数，将返回值改为封装的UTXO结构
3.修改getBalance函数调用