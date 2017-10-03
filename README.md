# github-pages
对原作者go语言爬取图片部分进行了补充
Author: yoghurtjia
Email：1527927373
1. 使用了并行机制，加快了图片的爬取速度
2. 将图片链接存放在redis数据库中
3. 使用go 1.9自带的sync.Map实现了图片url去重
