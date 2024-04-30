package structs_test

import "fmt"

// ? DISCOS extension .dsk

// Master Boot Record (MBR)
type MBR struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion [10]byte
	Mbr_dsk_signature  int32
	Dsk_fit            [1]byte
	Mbr_particion      [4]Partition
}

func PrintMBR(data MBR) {
	fmt.Printf("CreationDate: %s, fit: %s, size: %d \n",
		string(data.Mbr_fecha_creacion[:]),
		string(data.Dsk_fit[:]),
		data.Mbr_tamano)

	for i := 0; i < 4; i++ {
		fmt.Printf("Partition %d, Name: %s, Tipo: %s, Start: %d, Size: %d Status %s Correlativo %d ID %s CORRELATIVE: %d \n",
			i,
			string(data.Mbr_particion[i].Part_name[:]),
			string(data.Mbr_particion[i].Part_type[:]),
			data.Mbr_particion[i].Part_start,
			data.Mbr_particion[i].Part_size,
			string(data.Mbr_particion[i].Part_status[:]),
			data.Mbr_particion[i].Part_correlative,
			string(data.Mbr_particion[i].Part_id[:]),
			data.Mbr_particion[i].Part_correlative,
		)
	}
}

// Partition
type Partition struct {
	Part_status      [1]byte
	Part_type        [1]byte
	Part_fit         [1]byte
	Part_start       int32
	Part_size        int32
	Part_name        [16]byte
	Part_correlative int32
	Part_id          [4]byte
}

func PrintPartition(data Partition) {
	fmt.Printf("Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s\n",
		string(data.Part_name[:]),
		string(data.Part_type[:]),
		data.Part_start,
		data.Part_size,
		string(data.Part_status[:]),
		string(data.Part_id[:]))
}

func GetPartition(data Partition) string {
	str := fmt.Sprintf("Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s\n",
		string(data.Part_name[:]),
		string(data.Part_type[:]),
		data.Part_start,
		data.Part_size,
		string(data.Part_status[:]),
		string(data.Part_id[:]))
	return str
}

// Extended Boot Record (EBR)
type EBR struct {
	Part_mount [1]byte
	Part_fit   [1]byte
	Part_start int32
	Part_s     int32
	Part_next  int32
	Part_name  [16]byte
}

func PrintEBR(data EBR) {
	fmt.Printf("MOUNT: %s FIT: %s START: %d SIZE: %d NEXT: %d NAME: %s \n",
		string(data.Part_mount[:]),
		string(data.Part_fit[:]),
		data.Part_start,
		data.Part_s,
		data.Part_next,
		string(data.Part_name[:]),
	)
}

func GetEBR(data EBR) string {
	str := fmt.Sprintf("MOUNT: %s FIT: %s START: %d SIZE: %d NEXT: %d NAME: %s \n",
		string(data.Part_mount[:]),
		string(data.Part_fit[:]),
		data.Part_start,
		data.Part_s,
		data.Part_next,
		string(data.Part_name[:]),
	)
	return str
}

// ? CARPETAS Y ARCHIVOS (EXT3|EXT2)
type Superblock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [17]byte
	S_umtime            [17]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_fist_ino          int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

type Inode struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [17]byte
	I_ctime [17]byte
	I_mtime [17]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

func PrintInode(data Inode) {
	fmt.Printf("UID: %d \nGID: %d \nSIZE: %d \nACTUAL DATE: %s \nCREATION TIME: %s \nMODIFY TIME: %s \nBLOCKS:%d \nTYPE:%s \nPERM:%s \n",
		int(data.I_uid),
		int(data.I_gid),
		int(data.I_size),
		data.I_atime[:],
		data.I_ctime[:],
		data.I_mtime[:],
		data.I_block[:],
		data.I_type[:],
		data.I_perm[:],
	)
}

type Fileblock struct {
	B_content [64]byte
}

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type Folderblock struct {
	B_content [4]Content
}

type Pointerblock struct {
	B_pointers [16]int32
}

type Content_J struct {
	Operation [10]byte
	Path      [100]byte
	Content   [100]byte
	Date      [17]byte
}

type Journaling struct {
	Size      int32
	Ultimo    int32
	Contenido [50]Content_J
}
