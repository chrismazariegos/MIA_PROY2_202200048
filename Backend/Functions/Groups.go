package functions_test

import (
	"P1/Global"
	"P1/Structs"
	"P1/Utilities"
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

var (
	session      = false
	usuario      = Global.UserInfo{}
	groupCounter = 1
	userCounter  = 1
	blockIndex   = 0
	searchIndex  = 0
	letra        = ""
	ID           = ""
)

//?                    ADMINISTRACION DE USUARIOS Y GRUPOS
/* -------------------------------------------------------------------------- */
/*                                COMANDO LOGIN                               */
/* -------------------------------------------------------------------------- */
func ProcessLOGIN(input string, user *string, pass *string, id *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user":
			*user = flagValue
		case "pass":
			*pass = flagValue
		case "id":
			*id = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func LOGIN(user *string, pass *string, id *string) {

	letra = string((*id)[0])
	ID = *id

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

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// fmt.Println("Bitmap de bloques del inodo1")
	// fmt.Println(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}
	//fmt.Println("Fileblock " + fmt.Sprint(searchIndex))
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	userFound := false
	for _, line := range lines {
		// Imprimir cada línea
		// fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) > 3 {
			//fmt.Println("items[2]->" + items[2])
			if *user == items[len(items)-2] {
				userFound = true
				usuario.Nombre = items[len(items)-2]
				usuario.ID = items[0]
				session = true
				break
			}
		}
	}

	if !userFound {
		searchIndex++
		if searchIndex <= blockIndex {
			LOGIN(user, pass, id)
		} else {
			fmt.Println("Error: no se encontro al usuario")
			searchIndex = 0
			return
		}
	} else {
		Global.PrintUser(usuario)
		searchIndex = 0
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO LOGOUT                               */
/* -------------------------------------------------------------------------- */
func ProcessLOGOUT() {
	if session {
		println("Se ha cerrado la sesion")
		session = false
		searchIndex = 0
		return
	}
	println("Error: no hay una sesion activa")
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKGRP                               */
/* -------------------------------------------------------------------------- */
func ProcessMKGRP(input string, name *string) {
	if usuario.Nombre == "root" {
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
			default:
				fmt.Println("Error: Flag not found: " + flagName)
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func MKGRP(name *string) {
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

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// fmt.Println("Bitmap de bloques del inodo1")
	// fmt.Println(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	/* -------------------------------------------------------------------------- */
	/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
	/* -------------------------------------------------------------------------- */
	for _, line := range lines {
		// Imprimir cada línea
		fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) == 3 {
			if *name == items[2] {
				println("Error: nombre repetido")
				return
			}
		}
	}

	/* -------------------------------------------------------------------------- */
	/*                          PARSEAMOS LA INFORMACION                          */
	/* -------------------------------------------------------------------------- */
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	nuevoGrupo := fmt.Sprintf("%d,G,%s\n", groupCounter, *name)
	newContent := currentContent + nuevoGrupo

	/* -------------------------------------------------------------------------- */
	/*                 CREAMOS MAS FILEBLOCKS PARA GUARDAR LA INFO                */
	/* -------------------------------------------------------------------------- */
	if len(newContent) > len(Fileblock.B_content) {
		if blockIndex > int(len(crrInode.I_block)) {
			fmt.Println("Error: no hay mas bloques disponibles")
			return
		}
		blockIndex++
		var NEWFileblock structs_test.Fileblock
		if err := utilities_test.WriteObject(file, &NEWFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
			fmt.Println("Error reading Fileblock:", err)
			return
		}

		/* -------------------------------------------------------------------------- */
		/*                     ACTUALIZAMOS LOS BLOQUES DEL INODO                     */
		/* -------------------------------------------------------------------------- */
		crrInode.I_block[blockIndex] = 1

		if err := utilities_test.WriteObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
			fmt.Println("Error writing Inode to disk:", err)
			return
		}

		MKGRP(name)
	} else {
		/* -------------------------------------------------------------------------- */
		/*                GUARDA LA INFORMACION EN EL FILEBLOCK ACTUAL                */
		/* -------------------------------------------------------------------------- */
		copy(Fileblock.B_content[:], newContent)

		if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
			fmt.Println("Error writing Fileblock to disk:", err)
			return
		}

		println("ACTUALIZACION")
		// Mostrar el contenido actualizado del Fileblock
		data := string(Fileblock.B_content[:])
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		/* -------------------------------------------------------------------------- */
		/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
		/* -------------------------------------------------------------------------- */
		for _, line := range lines {
			// Imprimir cada línea
			fmt.Println(line)
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO RMGRP                               */
/* -------------------------------------------------------------------------- */
func ProcessRMGRP(input string, name *string) {
	if usuario.Nombre == "root" {
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
			default:
				fmt.Println("Error: Flag not found: " + flagName)
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func RMGRP(name *string) {
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

	// Leer el MBR del disco
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                               CARGAMOS EL MBR                              */
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
	/*                   LEEMOS EN INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                      LEEMOS EL CONTENIDO DEL FILEBLOCK                     */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                      COLOCAMOS EL STATUS DE ELIMINADO                      */
	/* -------------------------------------------------------------------------- */
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	lines := strings.Split(currentContent, "\n")
	deleted := false
	for i, line := range lines {
		if strings.Contains(line, *name) {
			lines[i] = "0,G," + *name
			deleted = true
			break
		}
	}

	/* -------------------------------------------------------------------------- */
	/*                   VERIFICAMOS BLOQUES O MENSAJE NOT FOUND                  */
	/* -------------------------------------------------------------------------- */
	if !deleted {
		searchIndex++
		if searchIndex > blockIndex {
			fmt.Println("Group not found")
			searchIndex = 0
			return
		}
		RMGRP(name)

	}

	/* -------------------------------------------------------------------------- */
	/*                          ACTUALIZAMOS EL CONTENIDO                         */
	/* -------------------------------------------------------------------------- */
	newContent := strings.Join(lines, "\n")
	copy(Fileblock.B_content[:], newContent)

	if deleted {
		if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
			fmt.Println("Error writing Fileblock to disk:", err)
			return
		}

		currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
		lines := strings.Split(currentContent, "\n")
		for i := range lines {
			println(lines[i])
		}

		searchIndex = 0
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKUSR                               */
/* -------------------------------------------------------------------------- */
func ProcessMKUSR(input string, user *string, pass *string, grp *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				*user = flagValue
			case "pass":
				*pass = flagValue
			case "grp":
				*grp = flagValue
			default:
				fmt.Println("Error: Flag not found: " + flagName)
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func MKUSR(user *string, pass *string, grp *string) {
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

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// fmt.Println("Bitmap de bloques del inodo1")
	// fmt.Println(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                          PARSEAMOS LA INFORMACION                          */
	/* -------------------------------------------------------------------------- */
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	searchIndex = 0
	var nuevoUsuario = BuscarGrupo(user, pass, grp)
	//fmt.Println("nuevo usuarios: " + nuevoUsuario)
	if nuevoUsuario == "" {
		fmt.Println("Error: No se encontro el grupo")
		return
	}
	newContent := currentContent + nuevoUsuario

	/* -------------------------------------------------------------------------- */
	/*                 CREAMOS MAS FILEBLOCKS PARA GUARDAR LA INFO                */
	/* -------------------------------------------------------------------------- */
	if len(newContent) > len(Fileblock.B_content) {
		if blockIndex > int(len(crrInode.I_block)) {
			fmt.Println("Error: no hay mas bloques disponibles")
			return
		}
		blockIndex++
		//fmt.Print("BlockIndex = " + fmt.Sprint(blockIndex))
		var NEWFileblock structs_test.Fileblock
		copy(NEWFileblock.B_content[:], nuevoUsuario)
		if err := utilities_test.WriteObject(file, &NEWFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
			fmt.Println("Error reading Fileblock:", err)
			return
		}

		println("MKUSR EXITOSO")
		// Mostrar el contenido actualizado del Fileblock
		data := string(NEWFileblock.B_content[:])
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		for _, line := range lines {
			// Imprimir cada línea
			fmt.Println(line)
		}

		/* -------------------------------------------------------------------------- */
		/*                     ACTUALIZAMOS LOS BLOQUES DEL INODO                     */
		/* -------------------------------------------------------------------------- */
		crrInode.I_block[blockIndex] = 1

		if err := utilities_test.WriteObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
			fmt.Println("Error writing Inode to disk:", err)
			return
		}
		searchIndex = 0

	} else {
		println("MKUSR EXITOSO")
		/* -------------------------------------------------------------------------- */
		/*                GUARDA LA INFORMACION EN EL FILEBLOCK ACTUAL                */
		/* -------------------------------------------------------------------------- */
		copy(Fileblock.B_content[:], newContent)

		if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
			fmt.Println("Error writing Fileblock to disk:", err)
			return
		}

		// Mostrar el contenido actualizado del Fileblock
		data := string(Fileblock.B_content[:])
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		/* -------------------------------------------------------------------------- */
		/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
		/* -------------------------------------------------------------------------- */
		for _, line := range lines {
			// Imprimir cada línea
			fmt.Println(line)
		}
		searchIndex = 0
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO RMUSR                               */
/* -------------------------------------------------------------------------- */
func ProcessRMUSR(input string, user *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				*user = flagValue
			default:
				fmt.Println("Error: Flag not found: " + flagName)
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func RMUSR(user *string) {
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

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

	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		items := strings.Split(line, ",")
		if len(items) > 3 {
			if *user == items[len(items)-2] {
				items[0] = "0" // Setear el ID a 0
				newLine := strings.Join(items, ",")
				copy(Fileblock.B_content[:], []byte(newLine))
				if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
					fmt.Println("Error writing Fileblock to disk:", err)
					return
				}
				return
			}
		}
	}

	searchIndex++
	if searchIndex <= blockIndex {
		RMUSR(user)
	} else {
		fmt.Println("User not found")
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO CHGRP                               */
/* -------------------------------------------------------------------------- */
func ProcessCHGRP(input string, user *string, grp *string) {
}

func CHGRP(user *string, grp *string) {
}

/* -------------------------------------------------------------------------- */
/*                                 AUXILIARES                                 */
/* -------------------------------------------------------------------------- */
func ImprimirBloques() {
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

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// fmt.Println("Bitmap de bloques del inodo1")
	// fmt.Println(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}
	fmt.Println("Fileblock " + fmt.Sprint(searchIndex))
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		// Imprimir cada línea
		fmt.Println(line)
	}

	if searchIndex < blockIndex {
		searchIndex++
		ImprimirBloques()
	} else {
		searchIndex = 0
	}
}

func BuscarGrupo(user *string, pass *string, grp *string) string {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return ""
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return ""
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
		return ""
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return ""
	}

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return ""
	}

	// fmt.Println("Bitmap de bloques del inodo1")
	// fmt.Println(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return ""
	}
	//fmt.Println("Fileblock " + fmt.Sprint(searchIndex))
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	groupFound := false
	var newUserLine string
	for _, line := range lines {
		// Imprimir cada línea
		//fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) == 3 {
			//fmt.Println("items[2]->" + items[2])
			if *grp == items[2] {
				groupFound = true
				newUserLine = fmt.Sprintf("%d,G,%s,%s,%s\n", userCounter, *grp, *user, *pass)
				userCounter++
				break
			}
		}
	}

	if !groupFound {
		searchIndex++
		if searchIndex <= blockIndex {
			return BuscarGrupo(user, pass, grp)
		}
	} else {
		return newUserLine
	}
	return ""
}
