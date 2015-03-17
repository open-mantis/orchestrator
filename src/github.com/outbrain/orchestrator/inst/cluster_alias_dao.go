/*
   Copyright 2015 Shlomi Noach, courtesy Booking.com

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package inst

import (
	"fmt"
	"github.com/outbrain/golib/log"
	"github.com/outbrain/golib/sqlutils"
	"github.com/outbrain/orchestrator/db"
)

// ReadClusterAliases reads the entrie cluster name aliases mapping
func ReadClusterAliases() error {
	query := fmt.Sprintf(`
		select 
			cluster_name,
			alias
		from 
			cluster_alias
		`)
	db, err := db.OpenOrchestrator()
	if err != nil {
		goto Cleanup
	}

	clusterAliasMap = make(map[string]string)
	err = sqlutils.QueryRowsMap(db, query, func(m sqlutils.RowMap) error {
		clusterAliasMap[m.GetString("cluster_name")] = m.GetString("alias")
		return err
	})
Cleanup:

	if err != nil {
		log.Errore(err)
	}
	return err

}

// WriteClusterAlias will write (and override) a single cluster name mapping
func WriteClusterAlias(clusterName string, alias string) error {
	writeFunc := func() error {
		db, err := db.OpenOrchestrator()
		if err != nil {
			return log.Errore(err)
		}

		_, err = sqlutils.Exec(db, `
			replace into  
					cluster_alias (cluster_name, alias)
				values
					(?, ?)
			`,
			clusterName,
			alias)
		if err != nil {
			return log.Errore(err)
		}

		return nil
	}
	return ExecDBWriteFunc(writeFunc)
}
