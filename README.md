# FreeFRP
通过fofa、hunter、shadon导出frp服务器地址，找出未设置密码的frp服务器（白嫖）

------

> 根据 [利用FOFA 白嫖上万台FRP服务器](https://www.t00ls.com/articles-67269.html) 这个文章用go重新写了一个，速度方面比python更快，设置200线程，1分50秒扫完1w台服务器，程序写的比较简陋，还请见谅

使用方法很简单，默认不需要设置什么参数，只需要把导出来的服务器ip和端口放在ip.txt里面，然后直接运行

fofa指纹：app="frp"

> Usage of ./frpscan:
>   -f string
>         扫描文件 (default "ip.txt")

>   -o string
>         输出文件 (default "success.txt")

>   -p string
>         代理，支持http和socks5

>   -t int
>         线程 (default 100)


![image-20230122122239116](https://raw.githubusercontent.com/Ggasdfg321/FreeFRP/main/image-20230122122239116.png)
