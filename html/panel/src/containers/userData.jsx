import React from 'react';
import errorHandler from "../utils/errorHandler";

let ud = {"ID":1,"CreatedAt":"2021-10-21T23:04:45.3150016+02:00","UpdatedAt":"2021-10-21T23:44:08.6927611+02:00","DeletedAt":null,"email":"nathanzy15@gmail.com","token":"","expire_at":"2021-11-20T23:44:08.6927611+01:00"}

function Spinner(){
    return <div className="spinner-border" role="status">
        <span className="sr-only">Loading...</span>
    </div>
}

export default function UserInfo() {
    function fetchUser(){
        fetch("/lazarus/user").then((res)=>res.json()).then(res=>setUser(res)).catch(errorHandler)
    }
    function update(){
        setUpdating(true)
        fetch("/lazarus/update").catch(errorHandler).finally(()=>{
            setUpdating(false)
            fetchUser()
        })
    }
    const [user,setUser] = React.useState(null)
    const [updating,setUpdating] = React.useState(false)
    React.useEffect(fetchUser,[])
    return (
        user?<div>
            <h3>Hello, {user.email}</h3>
            <p>过期时间: {user.expire_at}</p>
            <div><button onClick={update} className={"btn btn-primary"}>再续一个月</button> {updating?<Spinner/>:<></>}</div>
        </div>:<>Loading...</>
    )
}