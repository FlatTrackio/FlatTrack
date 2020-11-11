(defn default
  [{:keys [render partials] :as config}]
  (let [{:keys [head]} partials]
    [:html
     (head config)
     [:body
      [:main.def-wrapper
       [:aside#sidebar.def-sidebar.unfocused
        (render :sitemap {:sort-by :firn-order})]
       [:article.content
        ;; [:div (render :toc)] ;; Optional; add a table of contents
        [:div (render :file)]]]]]))
