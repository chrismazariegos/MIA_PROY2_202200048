import React, { useState } from 'react';
import CodeMirror from '@uiw/react-codemirror'

const Consola = () => {
    const [texto, setTexto] = useState('');
    const [comando, setComando] = useState('');

    const handleChange = (event) => {
        setTexto(event.target.value);
    };

    const almacenarTexto = async () => {
        try {
            const response = await fetch("http://3.147.67.57:3000/comandos", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: new URLSearchParams({
                    texto: texto,
                }),
            });
            const data = await response.text();
            setComando(data);
            setTexto('')
        } catch (error) {

            console.log(error);
        }
    };

    const handleFileChange = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function (event) {
                alert(event.target.result);

                // Una vez que el archivo es leído, envía su contenido al servidor
                fetch("http://3.147.67.57:3000/cargarArchivo", {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams({
                        fileContent: event.target.result
                    })
                })
                    .then(response => response.text())
                    .then(data => {
                        setComando(data); // Respuesta del servidor
                    })
            };
            reader.readAsText(file);
        }
    };

    return (
        <>
            <div className="container mt-5">
                <div className="row">
                    <div className="col-md-6 offset-md-3">
                        <h2 className="mb-4">Ingrese un comando</h2>
                        <div className="input-group mb-3">
                            <input type="text" value={texto} onChange={handleChange} className="form-control" placeholder="Ingrese un comando" />
                            <button className="btn btn-primary" onClick={almacenarTexto}>Guardar</button>
                            <input type="file" accept=".mia" onChange={handleFileChange} className="form-control" style={{ display: 'none' }} id="fileInput" />
                            <label htmlFor="fileInput" className="btn btn-secondary">Cargar archivo</label>
                        </div>
                    </div>
                </div>
            </div>
            <CodeMirror
                width='100%'
                height='60vh'
                readOnly='true'
                value={comando}
            />
        </>
    );
};

export default Consola;
