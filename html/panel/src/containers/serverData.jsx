import React from 'react';
import errorHandler from "../utils/errorHandler";
export default function ServerData(){
    function fetchServers(){
        fetch("/lazarus/s/json").then().catch(errorHandler)
    }
    React.useEffect(()=>{
        fetchServers()
        setTimeout(()=>{
            setServers([1,2,3])
        },2000)
    },[])
    const [servers,setServers] = React.useState([])
    return servers.length?<div>
        <ul className="list-group list-group-horizontal">
            <li className="list-group-item">Host</li>
            <li className="list-group-item">CPU</li>
            <li className="list-group-item">Mem</li>
            <li className="list-group-item">Active Connections</li>
            <li className="list-group-item">Data Quota</li>
        </ul>
        {servers.map(value => {
            return (
                <ul className="list-group list-group-horizontal">
                    <li className="list-group-item">An item</li>
                    <li className="list-group-item">A second item</li>
                    <li className="list-group-item">A third item</li>
                    <li className="list-group-item">A third item</li>
                    <li className="list-group-item">A third item</li>
                </ul>
            )
        })}
    </div>:<>Loading...</>
}