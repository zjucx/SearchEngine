package index

import (
  ""
)

type buffer struct {
	buf *int		/* 输入缓冲区 */
	length int		/* 缓冲区当前有多少个数 */
	offset int	/* 缓冲区读到了文件的哪个位置 */
	idx int		/* 缓冲区的指针 */
}
