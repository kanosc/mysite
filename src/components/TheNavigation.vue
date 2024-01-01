<template>
  <nav class="navbar navbar-expand-md navbar-dark bg-dark" aria-label="Fourth navbar example" @touchmove.prevent>
    <div class="container-fluid">
      <AppLink class="navbar-brand" to="/">Home</AppLink>

      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarsExample04"
        aria-controls="navbarsExample04" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>

      <div class="collapse navbar-collapse" id="navbarsExample04">
        <ul class="navbar-nav me-auto mb-2 mb-md-0">
          <li class="nav-item">
            <AppLink class="nav-link" aria-current="page" to="/">主页</AppLink>
          </li>
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" data-bs-toggle="dropdown" aria-expanded="false">月の博客</a>
            <ul class="dropdown-menu">
              <li>
                <AppLink class="dropdown-item" to="http://glassysky.cn:8080">Guest</AppLink>
              </li>
              <li>
                <AppLink class="dropdown-item" to="http://glassysky.cn:8080/wp-admin/index.php">Admin</AppLink>
              </li>
            </ul>
          </li>
          <li class="nav-item">
            <AppLink id="luck" class="nav-link" to="/fortune">随机壁纸</AppLink>
          </li>
          <li class="nav-item">
            <AppLink id="yunpan" class="nav-link" to="/filesys">文件系统</AppLink>
          </li>
          <li class="nav-item">
            <AppLink class="nav-link" to="/chatrooms">聊天室</AppLink>
          </li>
        </ul>
        <form role="search">
          <input class="form-control" type="search" placeholder="Search" aria-label="Search">
        </form>
        <div class="m-2">
          <button class="btn btn-danger"><a class="text-white" href="/logout" style="text-decoration:none">退出</a></button>
        </div>
      </div>
    </div>
  </nav>
  <!-- <div style="height: 50px;"></div> -->
</template>