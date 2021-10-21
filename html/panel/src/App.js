import UserInfo from "./containers/userData";
import ServerData from "./containers/serverData";

function App() {
    function logout(){
        fetch("/lazarus/logout").then(()=>{
            setTimeout(()=>{
                window.location.href = "/"
            },100)
        }).catch(err=>{console.log(err);alert("failed")})
    }
    return (
        <div className="container">
            <h2 style={{marginTop:10}}>
                Lazarus Panel
                <span style={{fontSize:5,float:'right',marginTop:10}} onClick={logout}>
                    <a>logout</a>
                </span>
            </h2>
            <hr/>
            <UserInfo/>
            <hr/>
            <ServerData/>
        </div>
    );
}

export default App;
