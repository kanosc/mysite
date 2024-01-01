<script setup>
import { ref } from 'vue'
const emit = defineEmits(['upload-success'])

const selectedFiles = ref(null)
const progressShow = ref(false)
const progressText = ref('0%')
const progreassWidth = ref(0)

const uploadFilesWithProgress = () => {
    let files = selectedFiles.value.files;
    if (files.length == 0) {
        alert('No file selected');
        return
    }

    // 创建XMLHttpRequest对象
    const xhr = new XMLHttpRequest();

    // 监听上传进度事件
    xhr.upload.addEventListener('progress', function (event) {
        if (event.lengthComputable) {
            let percent = Math.round((event.loaded / event.total) * 100);
            progreassWidth.value = percent
            progressText.value = percent + '%'
        }
    });

    // 监听上传完成事件
    xhr.addEventListener('load', function () {
        if (xhr.status !== 200) {
            alert(xhr.responseText)
            return
        }
        progreassWidth.value = 100
        progressText.value = '100%'
        alert('Uploaded success')
        progressShow.value = false
        emit('upload-success')
        // getFiles()


    });

    // 监听上传错误事件
    xhr.addEventListener('error', function () {
        alert(responseError)
        progressShow.value = false
        var responseError = xhr.responseText


    });

    // 开始上传
    xhr.open('POST', '/api/file/v1/testdir')
    var formData = new FormData()
    for (var i = 0; i < files.length; i++) {
        formData.append('upload', files[i])
    }
    progressShow.value = true
    xhr.send(formData)

}

const uploadFiles = async () => {
    let formData = new FormData()
    for (let i = 0; i < selectedFiles.value.files.length; i++) {
        formData.append('upload', selectedFiles.value.files[i])
    }
    const rsp = await fetch("/api/file/v1/testdir", {
        method: "POST",
        body: formData,
    })
    if (rsp.status === 200) {
        alert("uploaded success")
        getFiles()
    } else {
        alert(rsp.text())
    }
}
</script>
<template>
    <div>
        <button class="btn btn-info mt-2 float-start" data-bs-toggle="modal" data-bs-target="#myModal">
            upload
        </button>
        <div class="modal" id="myModal">
            <div class="modal-dialog">
                <div class="modal-content">
                    <!-- 模态框头部 -->
                    <div class="modal-header">
                        <span class="modal-title">Please select files</span>
                        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                    </div>
                    <!-- 模态框内容 -->
                    <div class="modal-body">
                        <div class="mb-3">
                            <form id="upload-form" enctype="multipart/form-data" action="/api/file/v1/testdir" method="POST"
                                @submit.prevent="uploadFilesWithProgress()">
                                <input class="form-control" type="file" name="upload" multiple ref="selectedFiles" />
                                <hr>
                                <button type="submit" class="btn btn-info">upload</button>
                            </form>
                        </div>

                        <div class="progress" v-if="progressShow">
                            <div class="progress-bar" :style="{ width: progreassWidth + '%' }">{{ progressText
                            }}
                            </div>
                        </div>
                    </div>
                    <!-- 模态框底部 -->
                    <div class="modal-footer">
                        <button type="button" class="btn btn-danger" data-bs-dismiss="modal">close</button>
                    </div>

                </div>
            </div>
        </div>
    </div>
</template>