function Submit(e) {
    var xhr = new XMLHttpRequest();
    var data = {};
    var post = {};
    var arr=[];
    data["author"] = document.getElementById("author").value;
    data["title"] = document.getElementById("topic").value;
    post["author"] = document.getElementById("author").value;
    post["text"] = document.getElementById("text").value;
    arr[0]=post;
    data["posts"] = arr;
    console.log(data);
    xhr.open("POST", "/topic", true);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
    var dataJson = JSON.stringify(data);
    xhr.send(dataJson);
    var pathArray = window.location.pathname.split( '/' );
    var newPathname = "";
    for (i = 0; i < pathArray.length-3; i++) {
        newPathname += "/";
        newPathname += pathArray[i];
    }
    window.location.href = window.location.protocol + "//" + window.location.host + "/" + newPathname;

}

function Cancel(e) {
    var pathArray = window.location.pathname.split( '/' );
    var newPathname = "";
    for (i = 0; i < pathArray.length-3; i++) {
        newPathname += "/";
        newPathname += pathArray[i];
    }
    window.location.href = window.location.protocol + "//" + window.location.host + "/" + newPathname;

}


//document.getElementById('myform').addEventListener("submit", OnSubmit);
