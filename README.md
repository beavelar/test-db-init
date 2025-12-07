# Test DB Init

Testing out running DB init scripts on Railway through their "Generate from `GitHub Repo`" option

Main thing of importance is that the build will run on each commit to `main` (or a different branch if a different branch is selected in Railway). Since it runs on each commit, it's important to set some environment variables or some sort of DB versioning with version checks to control when the scripts run or which scripts run if the DB should be versioned
