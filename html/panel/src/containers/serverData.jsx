import React from 'react';
import errorHandler from "../utils/errorHandler";
export default function ServerData(){
    function fetchServers(){
        fetch("/lazarus/s/json").then(res=>res.json()).then(res=>{
            setServers(res)
        }).catch(errorHandler)
    }
    React.useEffect(()=>{
        fetchServers()
    },[])
    const [servers,setServers] = React.useState({})
    return Object.keys(servers).length?
        <table className="table">
            <thead>
            <tr>
                <th scope="col">host</th>
                <th scope="col">cpu%</th>
                <th scope="col">mem%</th>
                <th scope="col">tcp connections</th>
                <th scope={"col"}>data</th>
            </tr>
            </thead>
            <tbody>
            {Object.values(servers).map(value => {
                return (
                    <tr>
                        <td>{value.host}</td>
                        <td>{value.cpu}</td>
                        <td>{value.mem}</td>
                        <td>{value.tcp}</td>
                        <td>{value.dataTotal}/{value.dataQuota}</td>
                    </tr>
                )
            })}
            </tbody>
        </table> :<>Loading...</>
}