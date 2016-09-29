namespace go tutorial

struct Node {
	1: i32 a,
	2: i32 b
}

struct SortDesc {
	1: required i32 limit,
	2: optional bool asc = true
}

service RpcService {
	i32 Plus(1: Node req);
	string Hello();
 	list<i32> Sort(1: SortDesc sd);
}
