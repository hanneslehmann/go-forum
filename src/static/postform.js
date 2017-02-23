function Submit(e) {
    var xhr = new XMLHttpRequest();
    id=document.getElementById("topicId").value;
    var data = {};
    data["author"] = document.getElementById("author").value;
    data["text"] = document.getElementById("text").value;
    xhr.open("POST", "/topic/"+id, true);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
    var dataJson = JSON.stringify(data);
    xhr.send(dataJson);
    var pathArray = window.location.pathname.split( '/' );
    var newPathname = "";
    for (i = 0; i < pathArray.length-1; i++) {
        newPathname += "/";
        newPathname += pathArray[i];
    }
    window.location.href = window.location.protocol + "//" + window.location.host + "/" + newPathname;

}

function Cancel(e) {
    var pathArray = window.location.pathname.split( '/' );
    var newPathname = "";
    for (i = 0; i < pathArray.length-1; i++) {
        newPathname += "/";
        newPathname += pathArray[i];
    }
    window.location.href = window.location.protocol + "//" + window.location.host + "/" + newPathname;

}


//document.getElementById('myform').addEventListener("submit", OnSubmit);
