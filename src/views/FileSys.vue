<script setup>
import { ref, watch } from 'vue'
import UploadFileButton from '../components/UploadFileButton.vue'
import { copyList } from '../common/common.js'

const fileData = ref('')
let originData = null
const getFiles = async () => {
    try {
        const res = await fetch("/api/file/v1/files?dirname=testdir")
        const rsp = await res.json()
        originData = rsp.FileList
        fileData.value = copyList(originData)
    } catch (e) {
        alert(e)
    }
}
getFiles()

const sortByName = () => {
    fileData.value.sort((v1, v2) => {
        if (v1.Name.toUpperCase() >= v2.Name.toUpperCase()) {
            return 1
        } else {
            return -1
        }
    })
    // fileData.value = originData.toSorted((v1, v2) => {
    //     if (v1.Name.toUpperCase() >= v2.Name.toUpperCase()) {
    //         return 1
    //     } else {
    //         return -1
    //     }
    // })
}

const sortByTime = () => {
    fileData.value.sort((v1, v2) => {
        if (v1.ModifyTime <= v2.ModifyTime) {
            return 1
        } else {
            return -1
        }
    })
    // fileData.value = originData.toSorted((v1, v2) => {
    //     if (v1.ModifyTime <= v2.ModifyTime) {
    //         return 1
    //     } else {
    //         return -1
    //     }
    // })
}

const searchKey = ref('')
watch(searchKey, (newKey, oldKey) => {
    if (newKey.trim().length > 0 && newKey !== '') {
        fileData.value = originData.filter(item => item.Name.indexOf(newKey) > -1)
    } else {
        fileData.value = copyList(originData)
    }
})

const deleteFile = async (fileName) => {
    if (confirm("Are you sure to delete " + fileName + " ?") == false) {
        return;
    }

    const response = await fetch("/api/file/v1/testdir/" + fileName, {
        method: "DELETE",
    })
    if (response.status === 200) {
        originData = originData.filter(item => item.Name !== fileName)
        fileData.value = copyList(originData)
    } else {
        alert(response.text)
    }
}
const makeDownloadURL = (filename) => { return "/api/file/v1/testdir/" + filename }

const copyLink = (fn) => {
    const completeLink = `${window.location.protocol}//${window.location.host}/api/file/v1/testdir/${fn}`

    // 创建一个输入框元素
    const inputElement = document.createElement('input');
    // 更新input的value
    inputElement.value = completeLink;
    // 将创建的input添加到容器中
    document.body.appendChild(inputElement);
    // 使用select方法将值选中
    inputElement.select();
    // 调用copy方法复制内容
    const flag = document.execCommand('copy');
    // 将输入框的值隐藏
    inputElement.style.display = 'none';
    // 根据返回值判断是否返回成功
    flag ? alert('address has been copied') : alert('copy address failed');
    // 删除生成的input元素
    inputElement.remove();


}
</script>
<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col-sm-2"></div>

            <div class="col-sm-8">
                <div class="mt-3 container-fluid">
                    <UploadFileButton @upload-success="getFiles()" />
                    <div class="btn-group float-start mx-2 mt-2" role="group">
                        <button type="button" class="btn btn-light dropdown-toggle" data-bs-toggle="dropdown"
                            aria-expanded="false">
                            Sort By
                        </button>
                        <ul class="dropdown-menu">
                            <li class="dropdown-item">
                                <button class="btn wide-btn" @click="sortByName()">
                                    Name
                                </button>
                            </li>
                            <li class="dropdown-item">
                                <button class="btn wide-btn" @click="sortByTime()">
                                    Time
                                </button>
                            </li>
                        </ul>
                    </div>
                    <form role="search" class="float-start mt-2">
                        <input class="form-control search-input" type="search" placeholder="Search" aria-label="Search"
                            v-model="searchKey">
                    </form>
                </div>
                <div class="table-responsive container-fluid">
                    <table class="table table-hover">
                        <thead>
                            <tr>
                                <th>File Name</th>
                                <th>Modifiy Time</th>
                                <th><span>Operation</span></th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="file in fileData">
                                <td>
                                    <a :href="makeDownloadURL(file.Name)" style="text-decoration: none; color: #32CD99;">{{
                                        file.Name }}</a>
                                </td>
                                <td>
                                    {{ file.ModifyTime }}
                                </td>
                                <td>
                                    <div class="btn-group">
                                        <button type="button" class="btn btn-primary dropdown-toggle btn-sm"
                                            data-bs-toggle="dropdown">handle</button>
                                        <ul class="dropdown-menu">
                                            <li class="dropdown-item">
                                                <a :href="makeDownloadURL(file.Name)" download
                                                    style="text-decoration: none; color: black;">download</a>
                                            </li>

                                            <li class="dropdown-item" @click="deleteFile(file.Name)">
                                                delete
                                            </li>

                                            <li class="dropdown-item" @click="copyLink(file.Name)">
                                                copy address
                                            </li>
                                        </ul>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <div class="col-sm-2"></div>
        </div>
    </div>
</template>
<style>
.btn-group li {
    cursor: pointer;
}
</style>
