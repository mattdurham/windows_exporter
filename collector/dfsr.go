// +build windows

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	registerCollector(
		"dfsr",
		func() collectorBuilder { return &DFRSConfig{} },
	)
}

// DFSRCollector contains the metric and state data of the DFSR collectors.
type DFSRCollector struct {
	// Meta
	dfsrScrapeDurationDesc *prometheus.Desc
	dfsrScrapeSuccessDesc  *prometheus.Desc

	// Connection source
	ConnectionBandwidthSavingsUsingDFSReplicationTotal *prometheus.Desc
	ConnectionBytesReceivedTotal                       *prometheus.Desc
	ConnectionCompressedSizeOfFilesReceivedTotal       *prometheus.Desc
	ConnectionFilesReceivedTotal                       *prometheus.Desc
	ConnectionRDCBytesReceivedTotal                    *prometheus.Desc
	ConnectionRDCCompressedSizeOfFilesReceivedTotal    *prometheus.Desc
	ConnectionRDCSizeOfFilesReceivedTotal              *prometheus.Desc
	ConnectionRDCNumberofFilesReceivedTotal            *prometheus.Desc
	ConnectionSizeOfFilesReceivedTotal                 *prometheus.Desc

	// Folder source
	FolderBandwidthSavingsUsingDFSReplicationTotal *prometheus.Desc
	FolderCompressedSizeOfFilesReceivedTotal       *prometheus.Desc
	FolderConflictBytesCleanedupTotal              *prometheus.Desc
	FolderConflictBytesGeneratedTotal              *prometheus.Desc
	FolderConflictFilesCleanedUpTotal              *prometheus.Desc
	FolderConflictFilesGeneratedTotal              *prometheus.Desc
	FolderConflictFolderCleanupsCompletedTotal     *prometheus.Desc
	FolderConflictSpaceInUse                       *prometheus.Desc
	FolderDeletedSpaceInUse                        *prometheus.Desc
	FolderDeletedBytesCleanedUpTotal               *prometheus.Desc
	FolderDeletedBytesGeneratedTotal               *prometheus.Desc
	FolderDeletedFilesCleanedUpTotal               *prometheus.Desc
	FolderDeletedFilesGeneratedTotal               *prometheus.Desc
	FolderFileInstallsRetriedTotal                 *prometheus.Desc
	FolderFileInstallsSucceededTotal               *prometheus.Desc
	FolderFilesReceivedTotal                       *prometheus.Desc
	FolderRDCBytesReceivedTotal                    *prometheus.Desc
	FolderRDCCompressedSizeOfFilesReceivedTotal    *prometheus.Desc
	FolderRDCNumberofFilesReceivedTotal            *prometheus.Desc
	FolderRDCSizeOfFilesReceivedTotal              *prometheus.Desc
	FolderSizeOfFilesReceivedTotal                 *prometheus.Desc
	FolderStagingSpaceInUse                        *prometheus.Desc
	FolderStagingBytesCleanedUpTotal               *prometheus.Desc
	FolderStagingBytesGeneratedTotal               *prometheus.Desc
	FolderStagingFilesCleanedUpTotal               *prometheus.Desc
	FolderStagingFilesGeneratedTotal               *prometheus.Desc
	FolderUpdatesDroppedTotal                      *prometheus.Desc

	// Volume source
	VolumeDatabaseLookupsTotal           *prometheus.Desc
	VolumeDatabaseCommitsTotal           *prometheus.Desc
	VolumeUSNJournalUnreadPercentage     *prometheus.Desc
	VolumeUSNJournalRecordsAcceptedTotal *prometheus.Desc
	VolumeUSNJournalRecordsReadTotal     *prometheus.Desc

	// Map of child collector functions used during collection
	dfsrChildCollectors []dfsrCollectorFunc
}

type DFRSConfig struct {
	dfsrEnabledCollectors string
}

func (config *DFRSConfig) RegisterFlags(application *kingpin.Application) {
	application.Flag("collectors.dfsr.sources-enabled", "Comma-seperated list of DFSR Perflib sources to use.").Default("connection,folder,volume").StringVar(&config.dfsrEnabledCollectors)
}

func (config *DFRSConfig) RegisterFlagsForLibrary(m map[string]string) {
	if dfrs, found := m["collectors.dfsr.sources-enabled"]; found == false {
		dfrs = "connection,folder,volume"
	} else {
		config.dfsrEnabledCollectors = dfrs

	}
}

func (config *DFRSConfig) Setup() {

}

func (config *DFRSConfig) Build() (Collector, error) {

	// Perflib sources are dynamic, depending on the enabled child collectors
	var perflibDependencies []string
	for _, source := range expandEnabledChildCollectors(config.dfsrEnabledCollectors) {
		perflibDependencies = append(perflibDependencies, dfsrGetPerfObjectName(source))
	}
	const subsystem = "dfsr"

	enabled := expandEnabledChildCollectors(config.dfsrEnabledCollectors)
	perfCounters := make([]string, 0, len(enabled))
	for _, c := range enabled {
		perfCounters = append(perfCounters, dfsrGetPerfObjectName(c))
	}
	addPerfCounterDependencies(subsystem, perfCounters)
	c := new(DFSRCollector)
	// meta
	c.dfsrScrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "collector_duration_seconds"),
		"windows_exporter: Duration of an dfsr child collection.",
		[]string{"collector"},
		nil,
	)
	c.dfsrScrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "collector_success"),
		"windows_exporter: Whether a dfsr child collector was successful.",
		[]string{"collector"},
		nil,
	)

	// Connection
	c.ConnectionBandwidthSavingsUsingDFSReplicationTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "bandwidth_savings_using_dfs_replication_bytes_total"),
		"Total amount of bandwidth savings using DFS Replication for this connection, in bytes",
		[]string{"name"},
		nil,
	)

	c.ConnectionBytesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "bytes_received_total"),
		"Total bytes received for connection",
		[]string{"name"},
		nil,
	)

	c.ConnectionCompressedSizeOfFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "compressed_size_of_files_received_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.ConnectionFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "received_files_total"),
		"Total number of files receieved for connection",
		[]string{"name"},
		nil,
	)

	c.ConnectionRDCBytesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "rdc_received_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.ConnectionRDCCompressedSizeOfFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "rdc_compressed_size_of_files_received_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.ConnectionRDCNumberofFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "rdc_received_files_total"),
		"Total number of Remote Differential Compression files received",
		[]string{"name"},
		nil,
	)

	c.ConnectionRDCSizeOfFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "rdc_size_of_received_files_bytes_total"),
		"Total size of received Remote Differential Compression files, in bytes.",
		[]string{"name"},
		nil,
	)

	c.ConnectionSizeOfFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "files_received_bytes_total"),
		"Total size of files received, in bytes",
		[]string{"name"},
		nil,
	)

	// Folder
	c.FolderBandwidthSavingsUsingDFSReplicationTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "bandwidth_savings_using_dfs_replication_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderCompressedSizeOfFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "compressed_size_of_files_received_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderConflictBytesCleanedupTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "conflict_cleaned_up_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderConflictBytesGeneratedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "conflict_generated_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderConflictFilesCleanedUpTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "conflict_cleaned_up_files_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderConflictFilesGeneratedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "conflict_generated_files_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderConflictFolderCleanupsCompletedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "conflict_folder_cleanups_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderConflictSpaceInUse = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "conflict_space_in_use_bytes"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderDeletedSpaceInUse = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "deleted_space_in_use_bytes"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderDeletedBytesCleanedUpTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "deleted_cleaned_up_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderDeletedBytesGeneratedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "deleted_generated_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderDeletedFilesCleanedUpTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "deleted_cleaned_up_files_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderDeletedFilesGeneratedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "deleted_generated_files_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderFileInstallsRetriedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "file_installs_retried_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderFileInstallsSucceededTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "file_installs_succeeded_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "received_files_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderRDCBytesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "rdc_received_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderRDCCompressedSizeOfFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "rdc_compressed_size_of_files_received_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderRDCNumberofFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "rdc_received_files_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderRDCSizeOfFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "rdc_files_received_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderSizeOfFilesReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "files_received_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderStagingSpaceInUse = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "staging_space_in_use_bytes"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderStagingBytesCleanedUpTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "staging_cleaned_up_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderStagingBytesGeneratedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "staging_generated_bytes_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderStagingFilesCleanedUpTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "staging_cleaned_up_files_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderStagingFilesGeneratedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "staging_generated_files_total"),
		"",
		[]string{"name"},
		nil,
	)

	c.FolderUpdatesDroppedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "dropped_updates_total"),
		"",
		[]string{"name"},
		nil,
	)

	// Volume
	c.VolumeDatabaseCommitsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "database_commits_total"),
		"Total number of DFSR Volume database commits",
		[]string{"name"},
		nil,
	)

	c.VolumeDatabaseLookupsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "database_lookups_total"),
		"Total number of DFSR Volume database lookups",
		[]string{"name"},
		nil,
	)

	c.VolumeUSNJournalUnreadPercentage = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "usn_journal_unread_percentage"),
		"Percentage of DFSR Volume USN journal records that are unread",
		[]string{"name"},
		nil,
	)

	c.VolumeUSNJournalRecordsAcceptedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "usn_journal_accepted_records_total"),
		"Total number of USN journal records accepted",
		[]string{"name"},
		nil,
	)

	c.VolumeUSNJournalRecordsReadTotal = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "usn_journal_read_records_total"),
		"Total number of DFSR Volume USN journal records read",
		[]string{"name"},
		nil,
	)

	c.dfsrChildCollectors = c.getDFSRChildCollectors(enabled)
	return c, nil
}

type dfsrCollectorFunc func(ctx *ScrapeContext, ch chan<- prometheus.Metric) error

// Map Perflib sources to DFSR collector names
// E.G. volume -> DFS Replication Service Volumes
func dfsrGetPerfObjectName(collector string) string {
	prefix := "DFS "
	suffix := ""
	switch collector {
	case "connection":
		suffix = "Replication Connections"
	case "folder":
		suffix = "Replicated Folders"
	case "volume":
		suffix = "Replication Service Volumes"
	}
	return (prefix + suffix)
}

// Maps enabled child collectors names to their relevant collection function,
// for use in DFSRCollector.Collect()
func (c *DFSRCollector) getDFSRChildCollectors(enabledCollectors []string) []dfsrCollectorFunc {
	var dfsrCollectors []dfsrCollectorFunc
	for _, collector := range enabledCollectors {
		switch collector {
		case "connection":
			dfsrCollectors = append(dfsrCollectors, c.collectConnection)
		case "folder":
			dfsrCollectors = append(dfsrCollectors, c.collectFolder)
		case "volume":
			dfsrCollectors = append(dfsrCollectors, c.collectVolume)
		}
	}

	return dfsrCollectors
}

// Collect implements the Collector interface.
// Sends metric values for each metric to the provided prometheus Metric channel.
func (c *DFSRCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	for _, fn := range c.dfsrChildCollectors {
		err := fn(ctx, ch)
		if err != nil {
			return err
		}
	}
	return nil
}

// Child-specific functions are provided to this function and executed concurrently.
// Child collector metrics & results are reported.
func (c *DFSRCollector) execute(ctx *ScrapeContext, name string, fn dfsrCollectorFunc, ch chan<- prometheus.Metric) error {
	// Child collector function called here, sends metric data back through channel
	err := fn(ctx, ch)
	if err != nil {
		return err
	}
	return err
}

// Perflib: "DFS Replication Service Connections"
type PerflibDFSRConnection struct {
	Name string

	BandwidthSavingsUsingDFSReplicationTotal float64 `perflib:"Bandwidth Savings Using DFS Replication"`
	BytesReceivedTotal                       float64 `perflib:"Total Bytes Received"`
	CompressedSizeOfFilesReceivedTotal       float64 `perflib:"Compressed Size of Files Received"`
	FilesReceivedTotal                       float64 `perflib:"Total Files Received"`
	RDCBytesReceivedTotal                    float64 `perflib:"RDC Bytes Received"`
	RDCCompressedSizeOfFilesReceivedTotal    float64 `perflib:"RDC Compressed Size of Files Received"`
	RDCNumberofFilesReceivedTotal            float64 `perflib:"RDC Number of Files Received"`
	RDCSizeOfFilesReceivedTotal              float64 `perflib:"RDC Size of Files Received"`
	SizeOfFilesReceivedTotal                 float64 `perflib:"Size of Files Received"`
}

func (c *DFSRCollector) collectConnection(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	var dst []PerflibDFSRConnection
	if err := unmarshalObject(ctx.perfObjects["DFS Replication Connections"], &dst); err != nil {
		return err
	}

	for _, connection := range dst {
		ch <- prometheus.MustNewConstMetric(
			c.ConnectionBandwidthSavingsUsingDFSReplicationTotal,
			prometheus.CounterValue,
			connection.BandwidthSavingsUsingDFSReplicationTotal,
			connection.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ConnectionBytesReceivedTotal,
			prometheus.CounterValue,
			connection.BytesReceivedTotal,
			connection.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ConnectionCompressedSizeOfFilesReceivedTotal,
			prometheus.CounterValue,
			connection.CompressedSizeOfFilesReceivedTotal,
			connection.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ConnectionFilesReceivedTotal,
			prometheus.CounterValue,
			connection.FilesReceivedTotal,
			connection.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ConnectionRDCBytesReceivedTotal,
			prometheus.CounterValue,
			connection.RDCBytesReceivedTotal,
			connection.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ConnectionRDCCompressedSizeOfFilesReceivedTotal,
			prometheus.CounterValue,
			connection.RDCCompressedSizeOfFilesReceivedTotal,
			connection.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ConnectionRDCSizeOfFilesReceivedTotal,
			prometheus.CounterValue,
			connection.RDCSizeOfFilesReceivedTotal,
			connection.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ConnectionRDCNumberofFilesReceivedTotal,
			prometheus.CounterValue,
			connection.RDCNumberofFilesReceivedTotal,
			connection.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ConnectionSizeOfFilesReceivedTotal,
			prometheus.CounterValue,
			connection.SizeOfFilesReceivedTotal,
			connection.Name,
		)

	}
	return nil
}

// Perflib: "DFS Replicated Folder"
type PerflibDFSRFolder struct {
	Name string

	BandwidthSavingsUsingDFSReplicationTotal float64 `perflib:"Bandwidth Savings Using DFS Replication"`
	CompressedSizeOfFilesReceivedTotal       float64 `perflib:"Compressed Size of Files Received"`
	ConflictBytesCleanedupTotal              float64 `perflib:"Conflict Bytes Cleaned Up"`
	ConflictBytesGeneratedTotal              float64 `perflib:"Conflict Bytes Generated"`
	ConflictFilesCleanedUpTotal              float64 `perflib:"Conflict Files Cleaned Up"`
	ConflictFilesGeneratedTotal              float64 `perflib:"Conflict Files Generated"`
	ConflictFolderCleanupsCompletedTotal     float64 `perflib:"Conflict Folder Cleanups Completed"`
	ConflictSpaceInUse                       float64 `perflib:"Conflict Space In Use"`
	DeletedSpaceInUse                        float64 `perflib:"Deleted Space In Use"`
	DeletedBytesCleanedUpTotal               float64 `perflib:"Deleted Bytes Cleaned Up"`
	DeletedBytesGeneratedTotal               float64 `perflib:"Deleted Bytes Generated"`
	DeletedFilesCleanedUpTotal               float64 `perflib:"Deleted Files Cleaned Up"`
	DeletedFilesGeneratedTotal               float64 `perflib:"Deleted Files Generated"`
	FileInstallsRetriedTotal                 float64 `perflib:"File Installs Retried"`
	FileInstallsSucceededTotal               float64 `perflib:"File Installs Succeeded"`
	FilesReceivedTotal                       float64 `perflib:"Total Files Received"`
	RDCBytesReceivedTotal                    float64 `perflib:"RDC Bytes Received"`
	RDCCompressedSizeOfFilesReceivedTotal    float64 `perflib:"RDC Compressed Size of Files Received"`
	RDCNumberofFilesReceivedTotal            float64 `perflib:"RDC Number of Files Received"`
	RDCSizeOfFilesReceivedTotal              float64 `perflib:"RDC Size of Files Received"`
	SizeOfFilesReceivedTotal                 float64 `perflib:"Size of Files Received"`
	StagingSpaceInUse                        float64 `perflib:"Staging Space In Use"`
	StagingBytesCleanedUpTotal               float64 `perflib:"Staging Bytes Cleaned Up"`
	StagingBytesGeneratedTotal               float64 `perflib:"Staging Bytes Generated"`
	StagingFilesCleanedUpTotal               float64 `perflib:"Staging Files Cleaned Up"`
	StagingFilesGeneratedTotal               float64 `perflib:"Staging Files Generated"`
	UpdatesDroppedTotal                      float64 `perflib:"Updates Dropped"`
}

func (c *DFSRCollector) collectFolder(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	var dst []PerflibDFSRFolder
	if err := unmarshalObject(ctx.perfObjects["DFS Replicated Folders"], &dst); err != nil {
		return err
	}

	for _, folder := range dst {
		ch <- prometheus.MustNewConstMetric(
			c.FolderBandwidthSavingsUsingDFSReplicationTotal,
			prometheus.CounterValue,
			folder.BandwidthSavingsUsingDFSReplicationTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderCompressedSizeOfFilesReceivedTotal,
			prometheus.CounterValue,
			folder.CompressedSizeOfFilesReceivedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderConflictBytesCleanedupTotal,
			prometheus.CounterValue,
			folder.ConflictBytesCleanedupTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderConflictBytesGeneratedTotal,
			prometheus.CounterValue,
			folder.ConflictBytesGeneratedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderConflictFilesCleanedUpTotal,
			prometheus.CounterValue,
			folder.ConflictFilesCleanedUpTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderConflictFilesGeneratedTotal,
			prometheus.CounterValue,
			folder.ConflictFilesGeneratedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderConflictFolderCleanupsCompletedTotal,
			prometheus.CounterValue,
			folder.ConflictFolderCleanupsCompletedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderConflictSpaceInUse,
			prometheus.GaugeValue,
			folder.ConflictSpaceInUse,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderDeletedSpaceInUse,
			prometheus.GaugeValue,
			folder.DeletedSpaceInUse,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderDeletedBytesCleanedUpTotal,
			prometheus.CounterValue,
			folder.DeletedBytesCleanedUpTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderDeletedBytesGeneratedTotal,
			prometheus.CounterValue,
			folder.DeletedBytesGeneratedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderDeletedFilesCleanedUpTotal,
			prometheus.CounterValue,
			folder.DeletedFilesCleanedUpTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderDeletedFilesGeneratedTotal,
			prometheus.CounterValue,
			folder.DeletedFilesGeneratedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderFileInstallsRetriedTotal,
			prometheus.CounterValue,
			folder.FileInstallsRetriedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderFileInstallsSucceededTotal,
			prometheus.CounterValue,
			folder.FileInstallsSucceededTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderFilesReceivedTotal,
			prometheus.CounterValue,
			folder.FilesReceivedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderRDCBytesReceivedTotal,
			prometheus.CounterValue,
			folder.RDCBytesReceivedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderRDCCompressedSizeOfFilesReceivedTotal,
			prometheus.CounterValue,
			folder.RDCCompressedSizeOfFilesReceivedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderRDCNumberofFilesReceivedTotal,
			prometheus.CounterValue,
			folder.RDCNumberofFilesReceivedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderRDCSizeOfFilesReceivedTotal,
			prometheus.CounterValue,
			folder.RDCSizeOfFilesReceivedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderSizeOfFilesReceivedTotal,
			prometheus.CounterValue,
			folder.SizeOfFilesReceivedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderStagingSpaceInUse,
			prometheus.GaugeValue,
			folder.StagingSpaceInUse,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderStagingBytesCleanedUpTotal,
			prometheus.CounterValue,
			folder.StagingBytesCleanedUpTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderStagingBytesGeneratedTotal,
			prometheus.CounterValue,
			folder.StagingBytesGeneratedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderStagingFilesCleanedUpTotal,
			prometheus.CounterValue,
			folder.StagingFilesCleanedUpTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderStagingFilesGeneratedTotal,
			prometheus.CounterValue,
			folder.StagingFilesGeneratedTotal,
			folder.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.FolderUpdatesDroppedTotal,
			prometheus.CounterValue,
			folder.UpdatesDroppedTotal,
			folder.Name,
		)
	}
	return nil
}

// Perflib: "DFS Replication Service Volumes"
type PerflibDFSRVolume struct {
	Name string

	DatabaseCommitsTotal           float64 `perflib:"Database Commits"`
	DatabaseLookupsTotal           float64 `perflib:"Database Lookups"`
	USNJournalRecordsReadTotal     float64 `perflib:"USN Journal Records Read"`
	USNJournalRecordsAcceptedTotal float64 `perflib:"USN Journal Records Accepted"`
	USNJournalUnreadPercentage     float64 `perflib:"USN Journal Records Unread Percentage"`
}

func (c *DFSRCollector) collectVolume(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	var dst []PerflibDFSRVolume
	if err := unmarshalObject(ctx.perfObjects["DFS Replication Service Volumes"], &dst); err != nil {
		return err
	}

	for _, volume := range dst {
		ch <- prometheus.MustNewConstMetric(
			c.VolumeDatabaseLookupsTotal,
			prometheus.CounterValue,
			volume.DatabaseLookupsTotal,
			volume.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.VolumeDatabaseCommitsTotal,
			prometheus.CounterValue,
			volume.DatabaseCommitsTotal,
			volume.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.VolumeUSNJournalRecordsAcceptedTotal,
			prometheus.CounterValue,
			volume.USNJournalRecordsAcceptedTotal,
			volume.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.VolumeUSNJournalRecordsReadTotal,
			prometheus.CounterValue,
			volume.USNJournalRecordsReadTotal,
			volume.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.VolumeUSNJournalUnreadPercentage,
			prometheus.GaugeValue,
			volume.USNJournalUnreadPercentage,
			volume.Name,
		)

	}
	return nil
}
