# 2 january 2015
# usage: awk -f icheck16.awk diskimage
# enter one inumber at a time
# TODO the usual
BEGIN {
	diskimage = ARGV[1]
	ARGC--
}
{
	for (i = 0; i < 16; i++) {
		n = $1 + i
		c = sprintf("debugfs -R 'ncheck %d' '%s' 2>/dev/null", n, diskimage)
		system(c)
	}
}
