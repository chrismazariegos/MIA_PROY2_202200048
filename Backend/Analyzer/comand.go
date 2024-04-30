package analyzer_test

import (
	functions_test "P1/Functions"
	utilities_test "P1/Utilities"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func Command(input string) {

	// Verificar si el input está vacío
	if input == "" {
		return // No hacer nada si el input está vacío
	}

	comando := input
	input = strings.ToLower(input)
	switch {
	case strings.HasPrefix(input, "mkdisk"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleMKDISKCommand(comando)

	case strings.HasPrefix(input, "rmdisk"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleRMDISKCommand(comando)

	case strings.HasPrefix(input, "fdisk"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleFDISKCommand(comando)

	case strings.HasPrefix(input, "mount"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleMOUNTCommand(comando)

	case strings.HasPrefix(input, "unmount"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleUNMOUNTCommand(comando)

	case strings.HasPrefix(input, "mkfs"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleMKFSCommand(comando)

	case strings.HasPrefix(input, "login"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleLOGINCommand(comando)

	case strings.HasPrefix(input, "logout"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleLOGOUTCommand()

	case strings.HasPrefix(input, "mkgrp"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleMKGRPCommand(comando)

	case strings.HasPrefix(input, "rmgrp"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleRMGRPCommand(comando)

	case strings.HasPrefix(input, "mkusr"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleMKUSRCommand(comando)

	case strings.HasPrefix(input, "rmusr"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleRMUSRCommand(comando)

	case strings.HasPrefix(input, "mkfile"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleMKFILECommand(comando)

	case strings.HasPrefix(input, "cat"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleCATCommand(comando)

	case strings.HasPrefix(input, "remove"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleREMOVECommand(comando)

	case strings.HasPrefix(input, "edit"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleEDITCommand(comando)

	case strings.HasPrefix(input, "rename"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleRENAMECommand(comando)

	case strings.HasPrefix(input, "mkdir"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleMKDIRCommand(comando)

	case strings.HasPrefix(input, "copy"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleCOPYCommand(comando)

	case strings.HasPrefix(input, "move"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleMOVECommand(comando)

	case strings.HasPrefix(input, "find"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleFINDCommand(comando)

	case strings.HasPrefix(input, "chown"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleCHOWNCommand(comando)

	case strings.HasPrefix(input, "chgrp"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleCHGRPCommand(comando)

	case strings.HasPrefix(input, "chmod"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleCHMODCommand(comando)

	case strings.HasPrefix(input, "pause"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handlePAUSECommand()

	case strings.HasPrefix(input, "execute"):
		handleEXECUTECommand(comando)

	case strings.HasPrefix(input, "rep"):
		fmt.Println(">>>>>>>>>>>>>>>>>>>>" + comando)
		handleREPCommand(comando)

	case strings.HasPrefix(input, "#"):
		//Ignora las sentencias del lado derecho
	case strings.HasPrefix(input, "exit"):
		fmt.Println("Saliendo del programa...")
		os.Exit(0)
	default:
		fmt.Println("Comando no reconocido:", comando)
	}
}

var (
	size        = flag.Int("size", 0, "Tamaño")
	fit         = flag.String("fit", "", "Ajuste")
	unit        = flag.String("unit", "", "Unidad")
	type_       = flag.String("type", "", "Tipo")
	driveletter = flag.String("driveletter", "", "Busqueda")
	name        = flag.String("name", "", "Nombre")
	delete      = flag.String("delete", "", "Eliminar")
	add         = flag.Int("add", 0, "Añadir/Quitar")
	path        = flag.String("path", "", "Directorio")
	id          = flag.String("id", "", "ID")
	fs          = flag.String("fs", "", "FDISK")
	ruta        = flag.String("ruta", "", "Ruta")
	user        = flag.String("user", "", "Usuario")
	pass        = flag.String("pass", "", "Password")
	grp         = flag.String("grp", "", "Group")
	r           = flag.Bool("r", false, "Rewrite")
	cont        = flag.String("cont", "", "Cont")
	destino     = flag.String("destino", "", "Destino")
	ugo         = flag.String("ugo", "", "UGO")
	file        = flag.String("file", "", "File to process")
	flagN       = flag.Bool("error", false, "Flag not found")
)

/* -------------------------------------------------------------------------- */
/*                           APLICACION DE COMANDOS                           */
/* -------------------------------------------------------------------------- */
func handleMKDISKCommand(input string) {

	flag.Parse()
	functions_test.ProcessMKDISK(input, size, fit, unit, flagN)

	if *flagN {
		*flagN = false
		return
	}

	// validate size > 0
	if *size <= 0 {
		utilities_test.AddText("Error: Size must be greater than 0")
		return
	}

	// validate fit equals to b/w/f
	if *fit != "b" && *fit != "f" && *fit != "w" {
		utilities_test.AddText("Error: Fit must be (bf/ff/wf)")
		return
	}

	// validate unit equals to k/m
	if *unit != "k" && *unit != "m" {
		utilities_test.AddText("Error: Unit must be (k/m)")
		return
	}

	//utilities_test.AddText("--------------------------------------------------------------------------")
	//utilities_test.AddText("                       	Discos creados correctamente                      ")
	//utilities_test.AddText("--------------------------------------------------------------------------")

	// Create the file
	functions_test.CreateBinFile(size, fit, unit)
	*size = 0
	*fit = ""
	*unit = ""
}

func handleRMDISKCommand(input string) {
	flag.Parse()
	functions_test.ProcessRMDISK(input, driveletter, flagN)

	if *flagN {
		*flagN = false
		return
	}

	// validate driveletter be a letter and not empty
	if !functions_test.ValidDriveLetter(*driveletter) {
		fmt.Println("Error: DriveLetter debe ser una letra")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error: DriveLetter es un campo obligatorio")
		return
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("                       MKDISK: PARAMETROS CORRECTOS                       ")
	fmt.Println("--------------------------------------------------------------------------")

	functions_test.DeleteBinFile(driveletter)
	*driveletter = ""
}

func handleFDISKCommand(input string) {
	flag.Parse()
	functions_test.ProcessFDISK(input, size, driveletter, name, unit, type_, fit, delete, add, path, flagN)

	if *flagN {
		*flagN = false
		return
	}

	//Obligatorio cuando no existe la particion
	// validate size > 0
	if *size <= 0 && *delete != "full" && *add == 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// validate driveletter be a letter and not empty
	if !functions_test.ValidDriveLetter(*driveletter) {
		fmt.Println("Error: DriveLetter must be a letter")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error: DriveLetter cannot be empty")
		return
	}

	// validate fit equals to b/w/f
	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error: Fit must be (BF/FF/WF)")
		return
	}

	// validate unit equals to b/k/m
	if *unit != "b" && *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be (B/K/M)")
		return
	}

	//println("ADD", *add)
	// validate type equals to P/E/L
	if *type_ != "p" && *type_ != "e" && *type_ != "l" && *delete != "full" && *add == 0 {
		fmt.Println("Error: Type must be (P/E/L)")
		return
	}

	if *delete != "" {
		if *delete != "full" {
			fmt.Println("Error: Delete must be full")
			return
		}
		if *name == "" && *path == "" {
			println("Error: you need path and name to delete")
			return
		}
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("                        FDISK: PARAMETROS CORRECTOS                       ")
	fmt.Println("--------------------------------------------------------------------------")

	functions_test.CRUD_Partitions(size, driveletter, name, unit, type_, fit, delete, add, path)
	*size = 0
	*driveletter = ""
	*name = ""
	*unit = ""
	*type_ = ""
	*fit = ""
	*delete = ""
	*add = 0
	*path = ""
}

func handleMOUNTCommand(input string) {
	flag.Parse()
	functions_test.ProcessMOUNT(input, driveletter, name, flagN)

	if *flagN {
		*flagN = false
		return
	}

	// validate driveletter be a letter and not empty
	if !functions_test.ValidDriveLetter(*driveletter) {
		fmt.Println("Error: DriveLetter must be a letter")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error: DriveLetter cannot be empty")
		return
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("                        MOUNT: PARAMETROS CORRECTOS                       ")
	fmt.Println("--------------------------------------------------------------------------")

	functions_test.MountPartition(driveletter, name)
	*driveletter = ""
	*name = ""
}

func handleUNMOUNTCommand(input string) {
	flag.Parse()
	functions_test.ProcessUNMOUNT(input, id, flagN)

	if *flagN {
		*flagN = false
		return
	}

	if *id == "" {
		println("Error: Id es un campo obligatorio")
	}

	letra := string((*id)[0])
	fmt.Println("DISCO:" + letra)

	if !functions_test.ValidDriveLetter(letra) {
		fmt.Println("Error: ID")
		fmt.Println("Error: DISCO INCORRECTO")
		return
	}

	numero := string((*id)[1])
	fmt.Println("PARTICION:" + numero)

	if !utilities_test.EsNumero(numero) {
		fmt.Println("Error: ID")
		fmt.Println("Error: PARTICION INCORRECTA")
		return
	}

	fmt.Println("CODIGO:" + string((*id)[2]) + string((*id)[3]))

	if string((*id)[2]) != "4" && string((*id)[3]) != "8" {
		fmt.Println("Error: ID")
		fmt.Println("Error: CODIGO INCORRECTO")
		return
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("                       UNMOUNT: PARAMETROS CORRECTOS                      ")
	fmt.Println("--------------------------------------------------------------------------")

	functions_test.UNMOUNT_Partition(id)
	*id = ""
}

func handleMKFSCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKFS(input, id, type_, fs, flagN)

	if *flagN {
		*flagN = false
		return
	}

	if *id == "" {
		println("Error: id es obligatorio")
	}

	if *fs != "2fs" && *fs != "3fs" {
		println("Error: fs debe ser 2fs o 3fs")
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("                         MKFS: PARAMETROS CORRECTOS                       ")
	fmt.Println("--------------------------------------------------------------------------")

	functions_test.MKFS(id, type_, fs)
	*id = ""
	*type_ = ""
	*fs = ""
}

/* -------------------------------------------------------------------------- */
/*                         ADMINISTRACION DE USUARIOS                         */
/* -------------------------------------------------------------------------- */
func handleLOGINCommand(input string) {
	flag.Parse()
	functions_test.ProcessLOGIN(input, user, pass, id)

	if *user == "" || *pass == "" || *id == "" {
		println("Error: campos incompletos")
	}

	functions_test.LOGIN(user, pass, id)

	*user = ""
	*pass = ""
	*id = ""
}

func handleLOGOUTCommand() {
	functions_test.ProcessLOGOUT()
}

func handleMKGRPCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKGRP(input, name)

	if *name == "" {
		println("Error: el campo name no puede estar vacio")
		return
	}

	functions_test.MKGRP(name)
	*name = ""
}

func handleRMGRPCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKGRP(input, name)

	if *name == "" {
		println("Error: el campo name no puede estar vacio")
		return
	}

	functions_test.RMGRP(name)
	*name = ""
}

func handleMKUSRCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKUSR(input, user, pass, grp)

	if len(*user) > 10 {
		println("Error: user no puede ser mayor a 10 caracteres")
		return
	}
	if len(*pass) > 10 {
		println("Error: password no puede ser mayor a 10 caracteres")
		return
	}
	if len(*grp) > 10 {
		println("Error: grupo no puede ser mayor a 10 caracteres")
		return
	}

	if *user == "" || *pass == "" || *grp == "" {
		println("Error: campos incompletos")
		return
	}

	functions_test.MKUSR(user, pass, grp)

	*user = ""
	*pass = ""
	*grp = ""
}

func handleRMUSRCommand(input string) {
	flag.Parse()
	functions_test.ProcessRMUSR(input, user)

	if *user == "" {
		println("Error: user no puede estar vacio")
		return
	}

	functions_test.RMUSR(user)

	*user = ""
}

func handleCHGRPCommand(input string) {
	flag.Parse()
	functions_test.ProcessCHGRP(input, user, grp)
}

/* -------------------------------------------------------------------------- */
/*                         ADMINISTRACION DE CARPETAS                         */
/* -------------------------------------------------------------------------- */
func handleMKDIRCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKDIR(input, path, r)

	if *path == "" {
		println("Error: path no puede estar vacio")
		return
	}

	fmt.Println("Path: " + *path)
	fmt.Print("r: ")
	fmt.Println(*r)

	functions_test.MKDIR(path, r)

	*path = ""
	*r = false
}

func handleMKFILECommand(input string) {
	flag.Parse()
	functions_test.ProcessMKFILE(input, path, r, size, cont)
}

func handleCATCommand(input string) {
	flag.Parse()
	functions_test.ProcessCAT(input, file)
}

func handleREMOVECommand(input string) {
	flag.Parse()
	functions_test.ProcessREMOVE(input, path)
}

func handleEDITCommand(input string) {
	flag.Parse()
	functions_test.ProcessEDIT(input, path, cont)
}

func handleRENAMECommand(input string) {
	flag.Parse()
	functions_test.ProcessRENAME(input, path, name)
}

func handleCOPYCommand(input string) {
	flag.Parse()
	functions_test.ProcessCOPY(input, path, destino)
}

func handleMOVECommand(input string) {
	flag.Parse()
	functions_test.ProcessMOVE(input, path, destino)
}

func handleFINDCommand(input string) {
	flag.Parse()
	functions_test.ProcessFIND(input, path, destino)
}

func handleCHOWNCommand(input string) {
	flag.Parse()
	functions_test.ProcessCHOWN(input, path, user, r)
}

func handleCHMODCommand(input string) {
	flag.Parse()
	functions_test.ProcessCHMOD(input, path, ugo, r)
}

/* -------------------------------------------------------------------------- */
/*                            COMANDOS AUXILIARES                             */
/* -------------------------------------------------------------------------- */
func handlePAUSECommand() {
	fmt.Println("Presione cualquier tecla para continuar...")
	fmt.Scanln() // Espera a que el usuario presione Enter
	fmt.Println("Continuando la ejecución...")
}

func handleEXECUTECommand(input string) {
	flag.Parse()
	functions_test.ProcessExecute(input, path, flagN)

	if *flagN {
		*flagN = false
		return
	}

	if *path == "" {
		fmt.Println("Error: Path cannot be empty")
		return
	}
	// Open bin file
	file, err := utilities_test.OpenFile(*path)
	if err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	// Crea un nuevo scanner para leer el archivo
	scanner := bufio.NewScanner(file)

	// Itera sobre cada línea del archivo
	for scanner.Scan() {
		linea := scanner.Text() // Lee la línea actual
		//fmt.Println(linea)
		Command(linea)
	}

	// Verifica si hubo algún error durante la lectura
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
	*path = ""
}

/* -------------------------------------------------------------------------- */
/*                                  REPORTES                                  */
/* -------------------------------------------------------------------------- */
func handleREPCommand(input string) {
	flag.Parse()
	functions_test.ProcessREP(input, name, path, id, ruta, flagN)

	if *flagN {
		*flagN = true
		return
	}

	if *name == "" || *path == "" || *id == "" {
		println("Error: incomplete statements")
		return
	}

	letra := string((*id)[0])
	fmt.Println("DISCO:" + letra)

	if !functions_test.ValidDriveLetter(letra) {
		fmt.Println("Error: ID")
		fmt.Println("Error: DISCO INCORRECTO")
		return
	}

	numero := string((*id)[1])
	fmt.Println("PARTICION:" + numero)

	if !utilities_test.EsNumero(numero) {
		fmt.Println("Error: ID")
		fmt.Println("Error: PARTICION INCORRECTA")
		return
	}

	fmt.Println("CODIGO:" + string((*id)[2]) + string((*id)[3]))

	if string((*id)[2]) != "4" && string((*id)[3]) != "8" {
		fmt.Println("Error: ID")
		fmt.Println("Error: CODIGO INCORRECTO")
		return
	}

	functions_test.GenerateReports(name, path, id, ruta)
}
