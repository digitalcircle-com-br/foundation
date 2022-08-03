package static

// func Run() error {
// 	core.Init("static")
// 	wd, _ := os.Getwd()
// 	logrus.Infof("Running from: %s", wd)
// 	files, err := os.ReadDir(wd)
// 	if err != nil {
// 		return err
// 	}

// 	for _, file := range files {
// 		logrus.Infof("	- %s", file.Name())
// 	}

// routemgr.HandleHttp("/_config", http.MethodGet, model.PERM_ALL, func(w http.ResponseWriter, r *http.Request) error {
// 	all := make(map[string]interface{})
// 	err := cfgmgr.Load("client", all)
// 	if err != nil {
// 		return err
// 	}
// 	var ret interface{}
// 	var ok bool

// 	k := r.URL.Query().Get("k")

// 	if k != "" {
// 		ret, ok = all[k]
// 		if !ok {
// 			http.Error(w, fmt.Sprintf("error retrieving key %s", k), http.StatusNotFound)
// 			return nil
// 		}
// 		json.NewEncoder(w).Encode(ret)
// 		return nil
// 	}
// 	host := strings.Split(r.Host, ":")[0]
// 	ret, ok = all[host]
// 	if !ok {
// 		ret, ok = all["default"]
// 		if !ok {
// 			http.Error(w, fmt.Sprintf("error retrieving host %s", host), http.StatusNotFound)
// 			return nil

// 		}
// 	}

// 	json.NewEncoder(w).Encode(ret)
// 	return nil

// })

// 	routemgr.Router().PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fname := strings.Split(r.URL.Path, "?")[0]

// 		if fname == "/" || fname == "" {
// 			fname = "/index.html"
// 		}
// 		fname = filepath.Join(wd, fname)
// 		logrus.Infof("Static - providing: %s", fname)
// 		mimetype := mime.TypeByExtension(filepath.Ext(fname))
// 		bs, err := os.ReadFile(fname)
// 		if err != nil {
// 			http.NotFound(w, r)
// 			return
// 		}
// 		w.Header().Add("content-type", mimetype)
// 		w.Write(bs)
// 	})
// 	return runmgr.RunS()
// }
