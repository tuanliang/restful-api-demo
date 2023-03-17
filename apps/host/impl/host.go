package impl

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/tuanliang/restful-api-demo/apps/host"
)

func (i *HostServiceImpl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	// 直接打印日志
	i.l.Debug("create host")
	// 带Format的日志打印，fmt.Sprintf()
	i.l.Debugf("create host %s ", ins.Name)
	// 携带额外的meta数据，常用于Trace系统
	i.l.With(logger.NewAny("request-id", "req01")).Debug("create host with meta kv")

	// 校验数据合法性
	if err := ins.Validate(); err != nil {
		fmt.Println("不合法", err)
		return nil, err
	}

	// 默认值填充
	ins.InjectDefault()

	if err := i.save(ctx, ins); err != nil {
		return nil, err
	}
	return ins, nil
}

func (i *HostServiceImpl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (*host.HostSet, error) {
	b := sqlbuilder.NewBuilder(QueryHostSQL)
	if req.Keywords != "" {
		b.Where("r.`name` LIKE ? OR r.description LIKE ? OR r.private_ip LIKE ? OR r.public_ip LIKE ?",
			"%"+req.Keywords+"%", "%"+req.Keywords+"%", req.Keywords+"%", req.Keywords+"%")
	}
	b.Limit(req.Offset(), req.GetPageSize())
	querySQL, args := b.Build()
	i.l.Debugf("query sql: %s, args: %v", querySQL, args)
	fmt.Println()
	stmt, err := i.db.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	set := host.NewHostSet()
	for rows.Next() {
		// 每扫描一行，就需要读取出来
		ins := host.NewHost()
		if err := rows.Scan(&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt,
			&ins.Type, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt,
			&ins.Account, &ins.PublicIP, &ins.PrivateIP, &ins.CPU, &ins.Memory, &ins.GPUSpec,
			&ins.GPUAmount, &ins.OSType, &ins.OSName, &ins.SerialNumber); err != nil {
			return nil, err
		}
		set.Add(ins)
	}
	countSql, args := b.BuildCount()
	i.l.Errorf("count sql: %s, args: %v\n", countSql, args)
	countStmt, err := i.db.PrepareContext(ctx, countSql)
	if err != nil {
		return nil, err
	}
	defer countStmt.Close()
	if err := countStmt.QueryRowContext(ctx, args...).Scan(&set.Total); err != nil {
		return nil, err
	}
	return set, nil
}

func (i *HostServiceImpl) DescribeHost(ctx context.Context, req *host.DescribeHostRequest) (*host.Host, error) {
	b := sqlbuilder.NewBuilder(QueryHostSQL)
	b.Where("r.id = ?", req.Id)

	querySQL, args := b.Build()
	i.l.Debugf("describe sql: %s, args: %v", querySQL, args)
	fmt.Println()
	stmt, err := i.db.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	ins := host.NewHost()
	err = stmt.QueryRowContext(ctx, args...).Scan(
		&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt,
		&ins.Type, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt,
		&ins.Account, &ins.PublicIP, &ins.PrivateIP, &ins.CPU, &ins.Memory, &ins.GPUSpec,
		&ins.GPUAmount, &ins.OSType, &ins.OSName, &ins.SerialNumber,
	)
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (i *HostServiceImpl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {
	// 获取已有对象
	ins, err := i.DescribeHost(ctx, host.NewDescribeHostRequestWithId(req.Id))
	if err != nil {
		return nil, err
	}
	// 根据跟新的模式，跟新对象
	switch req.UpdateMode {
	case host.UPDATE_MODE_PUT:
		if err := ins.Put(req.Host); err != nil {
			return nil, err
		}
	case host.UPDATE_MODE_PATCH:
		if err := ins.Patch(req.Host); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("update_mode only required put/patch")
	}
	// 检查跟新后的数据是否合法
	if err := ins.Validate(); err != nil {
		return nil, err
	}
	// 跟新数据库里面的数据
	if err := i.update(ctx, ins); err != nil {
		return nil, err
	}
	// 返回更新后的对象
	return ins, nil
}

func (i *HostServiceImpl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Host, error) {
	return nil, nil
}
