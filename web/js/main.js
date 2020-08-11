$(document).ready(function () {
    Main()
})

function Main() {
    UploadFile()
}

function UploadFile() {
    $(document).on("change",".upload-file",function () {
        let objFile = $(this)[0].files[0];
        let formData = new FormData();
        formData.append('file', objFile);
        $.ajax({
            url: '/upload',
            type: 'POST',
            cache: false,
            data: formData,
            processData: false,
            contentType: false,
            success: function (data) {
                if (data=="error"){
                    alert("上传错误")
                }else {
                    $.get("/excel/read",{name:data},function (result) {
                        if (result.length>0){
                            var temp=""
                            result[0].forEach(v=>{
                                temp+=`<th>${v}</th>`
                            })
                            var title=`<tr>${temp}</tr>`
                            $(".table-title").html(title)
                            var body=""
                            result.forEach((v,index)=>{
                                if (index!=0){
                                    var bodyItem=""
                                    v.forEach(vv=>{
                                        bodyItem+=`<td>${vv}</td>`
                                    })
                                    body+=`<tr>${bodyItem}</tr>`
                                }
                            })
                            $(".table-content").html(body)
                        }
                    })
                }
            }
        });

    })

}