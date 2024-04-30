package functions_test

import (
	structs_test "P1/Structs"
	utilities_test "P1/Utilities"
	"fmt"
	"regexp"
	"strings"
)

//?               ADMINISTRACION DE CARPETAS, ARCHIVOS Y PERMISOS
/* -------------------------------------------------------------------------- */
/*                                COMANDO MKDIR                               */
/* -------------------------------------------------------------------------- */
func ProcessMKDIR(input string, path *string, r *bool) {
	flags := strings.Split(input, "-")
	for _, i := range flags {
		if i == "r" {
			*r = true
		}
		f := strings.Split(i, "=")
		if f[0] == "path" {
			*path = f[1]
			if strings.Contains(f[1], " ") {
				*path = `"` + f[1] + `"`
			}
		}
	}

	re := regexp.MustCompile(`-(\w+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]

		switch flagName {
		case "r":
			*r = true
		case "path":
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func MKDIR(path *string, r *bool) {
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
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO MKFILE                               */
/* -------------------------------------------------------------------------- */
func ProcessMKFILE(input string, path *string, r *bool, size *int, cont *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "r":
			*r = true
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			*size = sizeValue
		case "cont":
			*cont = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func MKFILE(path *string, r *bool) {
}

/* -------------------------------------------------------------------------- */
/*                                 COMANDO CAT                                */
/* -------------------------------------------------------------------------- */
func ProcessCAT(input string, file *string) string {
	re := regexp.MustCompile(`-file(\d*)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagIndex := match[1]
		flagValue := match[2]

		// Eliminar comillas si están presentes en el valor
		flagValue = strings.Trim(flagValue, "\"")

		// Generar el nombre de la clave para el mapa

		// Asignar el valor al mapa
		*file = flagValue
		return flagIndex
	}
	return ""
}

func CAT(file *string) {
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO REMOVE                               */
/* -------------------------------------------------------------------------- */
func ProcessREMOVE(input string, path *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func REMOVE(path *string) {
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO EDIT                                */
/* -------------------------------------------------------------------------- */
func ProcessEDIT(input string, path *string, cont *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "cont":
			*cont = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func EDIT(path *string, cont *string) {
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO RENAME                               */
/* -------------------------------------------------------------------------- */
func ProcessRENAME(input string, path *string, name *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "name":
			*name = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func RENAME(path *string, name *string) {
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO COPY                                */
/* -------------------------------------------------------------------------- */
func ProcessCOPY(input string, path *string, destino *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "destino":
			*destino = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func COPY(path *string, destino *string) {
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MOVE                                */
/* -------------------------------------------------------------------------- */
func ProcessMOVE(input string, path *string, destino *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "destino":
			*destino = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func MOVE(path *string, destino *string) {
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO FIND                                */
/* -------------------------------------------------------------------------- */
func ProcessFIND(input string, path *string, destino *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "destino":
			*destino = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func FIND(path *string, destino *string) {
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO CHOWN                               */
/* -------------------------------------------------------------------------- */
func ProcessCHOWN(input string, path *string, user *string, r *bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "user":
			*user = flagValue
		case "r":
			*r = true
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func CHOWN(path *string, user *string, r *bool) {
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO CHMOD                               */
/* -------------------------------------------------------------------------- */
func ProcessCHMOD(input string, path *string, ugo *string, r *bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "ugo":
			*ugo = flagValue
		case "r":
			*r = true
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func CHMOD(path *string, ugo *string, r *bool) {
}
