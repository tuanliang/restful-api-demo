package impl

const (
	InsertResourceSQL = `
	INSERT INTO resource (
		id,
		vendor,
		region,
		create_at,
		expire_at,
		type,
		name,
		description,
		status,
		update_at,
		sync_at,
		accout,
		public_ip,
		private_ip
	)
	VALUES
		(?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`

	// INSERT INTO `host` ( resource_id, cpu, memory, gpu_amount, gpu_spec, os_type, os_name, serial_number )
	// VALUES
	// ( "111", 1, 2048, 1, 'n', 'linux', 'centos8', '00000' );
	InsertDescribeSQL = `
	INSERT INTO host ( resource_id, cpu, memory, gpu_amount, gpu_spec, os_type, os_name, serial_number )
	VALUES
		( ?,?,?,?,?,?,?,? );
	`

	queryHostSQL = `SELECT * FROM resource as r LEFT JOIN host h ON r.id=h.resource_id`
)
