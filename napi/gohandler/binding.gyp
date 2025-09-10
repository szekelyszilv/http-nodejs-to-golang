{
  "targets": [
    {
      "target_name": "gohandler",
      "sources": [ "src/addon.cc" ],
      "include_dirs": [
        "<!(node -p \"require('node-addon-api').include\")",
        "<(node_root_dir)/src"
      ],
      "dependencies": [
        "<!(node -p \"require('node-addon-api').targets\"):node_addon_api",
      ],
      "libraries": [
        "<(module_root_dir)/src/golib/gohandler.a"
      ],
    }
  ]
}
