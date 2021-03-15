package overpassapi

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/engelsjk/osm"
	"github.com/engelsjk/osm/osmapi"
)

type AugmentedDiffSeqNum uint64

func AugmentedDiffStatus(ctx context.Context) (AugmentedDiffSeqNum, error) {
	return DefaultDatasource.AugmentedDiffStatus(ctx)
}

func (ds *Datasource) AugmentedDiffStatus(ctx context.Context) (AugmentedDiffSeqNum, error) {

	url := ds.BaseURL + "/augmented_diff_status"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := ds.Client.Do(req.WithContext(ctx))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, &osmapi.UnexpectedStatusCodeError{
			Code: resp.StatusCode,
			URL:  url,
		}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	b := bytes.Trim(data, "\n")

	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0, err
	}

	return AugmentedDiffSeqNum(n), nil
}

// Changesets returns the complete list of changesets in for the given replication sequence.
func (ds *Datasource) AugmentedDiff(ctx context.Context, n AugmentedDiffSeqNum) (*osm.Diff, error) {
	url := fmt.Sprintf("%s/augmented_diff?id=%d&info=no", ds.BaseURL, n)

	diff := &osm.Diff{}
	if err := ds.getFromAPI(ctx, url, &diff); err != nil {
		return nil, err
	}

	return diff, nil
}
