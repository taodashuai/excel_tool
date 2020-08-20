$(document).ready(function () {
    Main()
})

function Main() {
    UploadFile()
}
let allName=[];
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
                        allName=[]
                        if (result.length>0){
                                result.forEach(value=>{
                                    allName.push(value["name"])
                                })
                        }
                        console.log(allName)
                    })
                }
            }
        });

    })

}