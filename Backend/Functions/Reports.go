package functions_test

import (
	structs_test "P1/Structs"
	utilities_test "P1/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// ?                     			REPORTES
func ProcessREP(input string, name *string, path *string, id *string, ruta *string, flagN *bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "name":
			*name = flagValue
		case "path":
			*path = flagValue
		case "id":
			*id = flagValue
		case "ruta":
			*ruta = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
			*flagN = true
		}
	}
}

func GenerateReports(name *string, path *string, id *string, ruta *string) {

	switch *name {
	case "mbr":
		REPORT_MBR(id, path)
	case "disk":
		REPORT_DISK(id, path)
	case "inode":
		REPORT_INODE(id, path)
	case "Journaling":
		REPORT_JOURNALING(id, path)
	case "block":
		REPORT_BLOCK(id, path)
	case "bm_inode":
		REPORT_BM_INODE(id, path)
	case "bm_block":
		REPORT_BM_BLOCK(id, path)
	case "tree":
		REPORT_TREE(path)
	case "sb":
		REPORT_SB(id, path)
	case "file":
		REPORT_FILE(id, path, ruta)
	case "ls":
		REPORT_LS(id, path, ruta)
	default:
		println("Reporte no reconocido:", *name)
	}
}

/* -------------------------------------------------------------------------- */
/*                               1 REPORTE MBR                                */
/* -------------------------------------------------------------------------- */
func REPORT_MBR(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)
	filepath := "./Disks/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var EPartition = false
	var EPartitionStart int

	var compareMBR structs_test.MBR
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")

	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	strP := ""
	strE := ""

	for _, partition := range TempMBR.Mbr_particion {
		partNameClean := strings.Trim(string(partition.Part_name[:]), "\x00")
		if partition.Part_correlative == 0 {
			continue
		} else {
			strP += fmt.Sprintf(`
		|Particion %d
		|{part_status|%s}
		|{part_type|%s}
		|{part_fit|%s}
		|{part_start|%d}
		|{part_size|%d}
		|{part_name|%s}`,
				partition.Part_correlative,
				string(partition.Part_status[:]),
				string(partition.Part_type[:]),
				string(partition.Part_fit[:]),
				partition.Part_start,
				partition.Part_size,
				partNameClean,
			)
		}

		//?EBR verificacion
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) && EPartition {
			// Validar que si no existe una particion extendida no se puede crear una logica
			//?EBR verificacion
			var x = 0
			for x < 1 {
				var TempEBR structs_test.EBR
				if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if EPartitionStart != 0 && TempEBR.Part_next != -1 {
					partNameClean := strings.Trim(string(TempEBR.Part_name[:]), "\x00")
					strE += fmt.Sprintf(`
		|Particion Logica
		|{part_status|%s}
		|{part_next|%d}
		|{part_fit|%s}
		|{part_start|%d}
		|{part_size|%d}
		|{part_name|%s}`,
						string(TempEBR.Part_mount[:]),
						TempEBR.Part_next,
						string(TempEBR.Part_fit[:]),
						TempEBR.Part_start,
						TempEBR.Part_s,
						partNameClean,
					)
					//print("fit logica")
					//println(string(TempEBR.Part_fit[:]))
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					//print("fit logica")
					//println(string(TempEBR.Part_fit[:]))
					partNameClean := strings.Trim(string(TempEBR.Part_name[:]), "\x00")
					strE += fmt.Sprintf(`
		|Particion Logica
		|{part_status|%s}
		|{part_next|%d}
		|{part_fit|%s}
		|{part_start|%d}
		|{part_size|%d}
		|{part_name|%s}`,
						string(TempEBR.Part_mount[:]),
						TempEBR.Part_next,
						string(TempEBR.Part_fit[:]),
						TempEBR.Part_start,
						TempEBR.Part_s,
						partNameClean,
					)
					strP += strE
					x = 1
				}
			}

		}

	}

	//structs_test.PrintMBR(TempMBR)

	dotCode := fmt.Sprintf(`
		digraph G {
 			fontname="Helvetica,Arial,sans-serif"
			node [fontname="Helvetica,Arial,sans-serif"]
			edge [fontname="Helvetica,Arial,sans-serif"]
			concentrate=True;
			rankdir=TB;
			node [shape=record];

			title [label="Reporte MBR" shape=plaintext fontname="Helvetica,Arial,sans-serif"];

  			mbr[label="
				{MBR: %s.dsk|
					{mbr_tamaño|%d}
					|{mbr_fecha_creacion|%s}
					|{mbr_disk_signature|%d}
								%s
				}
			"];
			title2 [label="Reporte EBR" shape=plaintext fontname="Helvetica,Arial,sans-serif"];
			
			ebr[label="
				{EBR%s}
			"];

			title -> mbr [style=invis];
    		mbr -> title2[style=invis];
			title2 -> ebr[style=invis];
		}`,
		letra,
		TempMBR.Mbr_tamano,
		TempMBR.Mbr_fecha_creacion,
		TempMBR.Mbr_dsk_signature,
		strP,
		strE,
	)

	// Escribir el contenido en el archivo DOT
	dotFilePath := "./Reports/mbr_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	// Llamar a Graphviz para generar el gráfico
	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("Reporte MBR, EBR generado en", pngFilePath)
}

/* -------------------------------------------------------------------------- */
/*                              2 REPORTE DISK                                */
/* -------------------------------------------------------------------------- */

func REPORT_DISK(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	filepath := "./Disks/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var EPartition = false
	var EPartitionStart int

	var compareMBR structs_test.MBR
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")

	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	strP := ""
	lastSize := int(TempMBR.Mbr_tamano)
	counter := -1
	for _, partition := range TempMBR.Mbr_particion {
		counter++
		if partition.Part_correlative == 0 {
			porcentaje := utilities_test.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			lastSize -= int(partition.Part_size)
			if porcentaje > 0 {
				strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
			}
		}

		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[0].Part_type[:]) {
			//println("primaria: " + string(partition.Part_name[:]))
			porcentaje := utilities_test.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			lastSize -= int(partition.Part_size)
			strP += fmt.Sprintf(`|Primaria\n%d%%`, porcentaje)
		}

		//?EBR verificacion
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) && EPartition {
			porcentaje := utilities_test.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			lastSize -= int(partition.Part_size)
			//println("extendida size")
			//println(partition.Part_size)
			strP += fmt.Sprintf(`|{Extendida %d%%|{`, porcentaje)
			// Validar que si no existe una particion extendida no se puede crear una logica
			//?EBR verificacion
			var x = 0
			for x < 1 {
				var TempEBR structs_test.EBR
				if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_next != -1 {
					if !bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|EBR|Particion logica %d%%`, porcentaje)
					} else {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						if porcentaje > 0 {
							strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
						}
					}
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					if !bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|EBR|Particion logica %d%%`, porcentaje)
					} else {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						if porcentaje > 0 {
							strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
						}
					}
					strP += "}}"
					x = 1
				}
			}
		}
	}
	porcentaje := utilities_test.CalcularPorcentaje(int64(lastSize), int64(TempMBR.Mbr_tamano))
	fmt.Print("PORCENTAJE RESTANTE: ")
	println(porcentaje)
	if porcentaje > 0 {
		strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
	}
	strP += "}"

	//structs_test.PrintMBR(TempMBR)

	dotCode := fmt.Sprintf(`
		digraph G {
 			fontname="Helvetica,Arial,sans-serif"
			node [fontname="Helvetica,Arial,sans-serif"]
			edge [fontname="Helvetica,Arial,sans-serif"]
			concentrate=True;
			rankdir=TB;
			node [shape=record];

			title [label="Reporte DISK %s" shape=plaintext fontname="Helvetica,Arial,sans-serif"];

  			dsk[label="
				{MBR}%s
				}
			"];
			
			title -> dsk [style=invis];
		}`,
		letra,
		strP,
	)

	// Escribir el contenido en el archivo DOT
	dotFilePath := "./Reports/disk_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	// Llamar a Graphviz para generar el gráfico
	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico de disco:", err)
		return
	}

	fmt.Println("Reporte DISK generado en", pngFilePath)
	println("MBR")
	structs_test.PrintMBR(TempMBR)
}

/* -------------------------------------------------------------------------- */
/*                              3 REPORTE INODE                               */
/* -------------------------------------------------------------------------- */

func REPORT_INODE(id *string, path *string) {
}

/* -------------------------------------------------------------------------- */
/*                              4 REPORTE BLOCK                               */
/* -------------------------------------------------------------------------- */

func REPORT_BLOCK(id *string, path *string) {

}

/* -------------------------------------------------------------------------- */
/*                            5 REPORTE BM_INODE                              */
/* -------------------------------------------------------------------------- */

func REPORT_BM_INODE(id *string, path *string) {

}

/* -------------------------------------------------------------------------- */
/*                             6 REPORTE BM_BLOC                              */
/* -------------------------------------------------------------------------- */

func REPORT_BM_BLOCK(id *string, path *string) {
}

/* -------------------------------------------------------------------------- */
/*                              7 REPORTE TREE                                */
/* -------------------------------------------------------------------------- */
func REPORT_TREE(path *string) {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	str := `digraph {
    node [shape=plaintext]
    rankdir = LR

    // Título del reporte
    title [label="Reporte TREE" shape=plaintext fontname="Helvetica,Arial,sans-serif" rank=max];
`

	/* -------------------------------------------------------------------------- */
	/*                       BUCLE PARA RECORRER LOS INODOS                       */
	/* -------------------------------------------------------------------------- */

	for i := 0; i < int(tempSuperblock.S_inodes_count); i++ {
		/* -------------------------------------------------------------------------- */
		/*                   			LEEMOS EL INODO                               */
		/* -------------------------------------------------------------------------- */
		indexInode := int32(i)
		var crrInode structs_test.Inode
		if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
			fmt.Println("Error reading inode:", err)
			return
		}

		str = str + fmt.Sprintf(`i%d [
        label=<
            <table border="1" cellborder="1" cellspacing="0">
                <tr><td colspan="2">Inodo %d</td></tr>
                <tr><td>UID</td><td>%d</td></tr>
                <tr><td>GIU</td><td>%d</td></tr>
                <tr><td>SIZE</td><td>%d</td></tr>
                <tr><td>A_TIME</td><td>%s</td></tr>
                <tr><td>C_TIME</td><td>%s</td></tr>
                <tr><td>M_TIME</td><td>%s</td></tr>
                <tr><td>BLOCK</td><td>%d</td></tr>
                <tr><td>TYPE</td><td>%s</td></tr>
                <tr><td>PERM</td><td>%s</td></tr>
            </table>
        >
    ]`,
			int(indexInode),
			int(indexInode),
			int(crrInode.I_uid),
			int(crrInode.I_gid),
			int(crrInode.I_size),
			crrInode.I_atime[:],
			crrInode.I_ctime[:],
			crrInode.I_mtime[:],
			crrInode.I_block[:],
			crrInode.I_type[:],
			crrInode.I_perm[:],
		)

		fmt.Print("--------------------------------------------------------------------------INODO ")
		fmt.Print(i)
		fmt.Println("-------------------------------------------------------------------------")
		structs_test.PrintInode(crrInode)
		var bloques int
		for j := 0; j < 15; j++ {
			if crrInode.I_block[j] == 1 {
				bloques += 1
			}
		}

		/* -------------------------------------------------------------------------- */
		/*                     BUCLE PARA MOSTRAR TODO LOS BLOQUES                    */
		/* -------------------------------------------------------------------------- */
		fmt.Println("\n\n-------------------BLOQUE---------------------")
		for j := 0; j < bloques; j++ {
			var Fileblock structs_test.Fileblock
			var Folderblock structs_test.Folderblock
			if i == 0 {
				searchIndex = 0
				if err := utilities_test.ReadObject(file, &Folderblock, int64(tempSuperblock.S_block_start)); err != nil {
					fmt.Println("Error reading Fileblock:", err)
					return
				}
				data := Folderblock.B_content[:]
				// Dividir la cadena en líneas

				/* -------------------------------------------------------------------------- */
				/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
				/* -------------------------------------------------------------------------- */

				str = str + `b0 [
        label=<
            <table border="1" cellborder="1" cellspacing="0">
			<tr><td colspan="2">Bloque 0</td></tr>`
				for _, line := range data {
					// Imprimir cada línea
					str = str + fmt.Sprintf(`
                <tr><td>%s</td><td>%d</td></tr>`,
						string(utilities_test.LimpiarCerosBinarios(line.B_name[:])),
						line.B_inodo,
					)
				}
				str = str + `</table>
        >
    ]
	
		i0 -> b0

		b0 -> i1
		`

			} else {
				if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(j))); err != nil {
					fmt.Println("Error reading Fileblock:", err)
					return
				}
				data := string(Fileblock.B_content[:])
				// Dividir la cadena en líneas
				lines := strings.Split(data, "\n")

				/* -------------------------------------------------------------------------- */
				/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
				/* -------------------------------------------------------------------------- */
				str = str + fmt.Sprintf(`b%d [
        label=<
            <table border="1" cellborder="1" cellspacing="0">
			<tr><td colspan="2">Bloque %d</td></tr>`, searchIndex, searchIndex)
				for z := 0; z < len(lines)-1; z++ {
					// Imprimir cada línea
					str = str + fmt.Sprintf(`
					<tr><td>%s</td></tr>`, lines[z])

				}
				str = str + `</table>
        >
    ]`

			}

			str = str + fmt.Sprintf(`
			i1 -> b%d
			`, searchIndex)
			searchIndex++
		}
	}
	str = str + "}"
	// Escribir el contenido en el archivo DOT
	dotFilePath := "./Reports/tree_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(str)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	// Llamar a Graphviz para generar el gráfico
	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("Reporte TREE generado en", pngFilePath)
}

/* -------------------------------------------------------------------------- */
/*                               8 REPORTE SB                                 */
/* -------------------------------------------------------------------------- */

func REPORT_SB(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                      GENERAMOS EL REPORTE EN GRAPHVIZ                      */
	/* -------------------------------------------------------------------------- */

	dotCode := fmt.Sprintf(`
		digraph G {
 			fontname="Helvetica,Arial,sans-serif"
			node [fontname="Helvetica,Arial,sans-serif"]
			edge [fontname="Helvetica,Arial,sans-serif"]
			concentrate=True;
			rankdir=TB;
			node [shape=record];

			title [label="Reporte SUPERBLOCK" shape=plaintext fontname="Helvetica,Arial,sans-serif"];

  			sb[label="
				{Superblock|
					{S_filesystem_type|%d}
					|{S_inodes_count|%d}
					|{S_blocks_count|%d}
					|{S_free_blocks_count|%d}
					|{S_free_inodes_count|%d}
					|{S_mtime|%s}
					|{S_umtime|%s}
					|{S_mnt_count|%d}
					|{S_magic|%d}
					|{S_inode_size|%d}
					|{S_block_size|%d}
					|{S_fist_ino|%d}
					|{S_first_blo|%d}
					|{S_bm_inode_start|%d}
					|{S_bm_block_start|%d}
					|{S_inode_start|%d}
					|{S_block_start|%d}
				}
			"];
			

			title -> sb [style=invis];
		}`,
		int(tempSuperblock.S_filesystem_type),
		int(tempSuperblock.S_inodes_count),
		int(tempSuperblock.S_blocks_count),
		int(tempSuperblock.S_free_blocks_count),
		int(tempSuperblock.S_free_inodes_count),
		tempSuperblock.S_mtime[:],
		tempSuperblock.S_umtime[:],
		int(tempSuperblock.S_mnt_count),
		int(tempSuperblock.S_magic),
		int(tempSuperblock.S_inode_size),
		int(tempSuperblock.S_block_size),
		int(tempSuperblock.S_fist_ino),
		int(tempSuperblock.S_first_blo),
		int(tempSuperblock.S_bm_inode_start),
		int(tempSuperblock.S_bm_block_start),
		int(tempSuperblock.S_inode_start),
		int(tempSuperblock.S_block_start),
	)

	// Escribir el contenido en el archivo DOT
	dotFilePath := "./Reports/sb_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	// Llamar a Graphviz para generar el gráfico
	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("Reporte SUPERBLOCK generado en", pngFilePath)
}

/* -------------------------------------------------------------------------- */
/*                              9 REPORTE FILE                                */
/* -------------------------------------------------------------------------- */

func REPORT_FILE(id *string, path *string, ruta *string) {
}

/* -------------------------------------------------------------------------- */
/*                              10 REPORTE LS                                 */
/* -------------------------------------------------------------------------- */

func REPORT_LS(id *string, path *string, ruta *string) {
}

/* -------------------------------------------------------------------------- */
/*                          11 REPORTE JOURNALING                             */
/* -------------------------------------------------------------------------- */

func REPORT_JOURNALING(id *string, path *string) {
}
