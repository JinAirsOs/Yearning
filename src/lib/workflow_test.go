package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLFormat(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "测试小于255sql",
			args: args{sql: "update mysql_cdc_es_config set es_action='update' where id = 2;"},
			want: "update mysql_cdc_es_config set es_action='update' where id = 2;",
		},
		{
			name: "测试大于255sql",
			args: args{sql: "CREATE TABLE `config_owner` (\n  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',\n  `owner_id` varchar(32) NOT NULL DEFAULT '' COMMENT '认领人id',\n  `owner_name` varchar(100) NOT NULL DEFAULT '' COMMENT '认领人名称',\n  `owner_grade` varchar(32) NOT NULL DEFAULT '' COMMENT '认领人等级：A B C D',\n  `owner_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 0:不在线 1:在线',\n  `is_deleted` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除，0：未删除，1：已删除',\n  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建记录时间',\n  `create_by` varchar(32) NOT NULL DEFAULT '' COMMENT '创建人',\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',\n  `update_by` varchar(32) NOT NULL DEFAULT '' COMMENT '最后修改人',\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `uk_owner_id` (`owner_id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='配置-认领人表';\n \n \n \nCREATE TABLE `issue_auto_owner` (\n  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',\n  `issue_id` int(11) unsigned NOT NULL DEFAULT 0 COMMENT 'issue id',\n  `owner_id` varchar(32) NOT NULL DEFAULT '' COMMENT '认领人id',\n  `is_deleted` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除，0：未删除，1：已删除',\n  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建记录时间',\n  `create_by` varchar(32) NOT NULL DEFAULT '' COMMENT '创建人',\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',\n  `update_by` varchar(32) NOT NULL DEFAULT '' COMMENT '最后修改人',\n  PRIMARY KEY (`id`),\n  KEY `idx_issue_id` (`issue_id`)\n) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='问题自动认领记录表';"},
			want: "CREATE TABLE `config_owner` ( \n  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键', \n  `owner_id` varchar(32) NOT NULL DEFAULT '' COMMENT '认领人id', \n  `owner_name` varchar(100) NOT NULL DEFAULT '' COMMENT '认领人名称', \n  `owner_grade` varchar(32) NOT NULL DEFAULT '' COMMENT '认领人等级：A B C D', \n  `owner_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 0:不在线 1:在线', \n  `is_deleted` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除，0：未删除，1：已删除', \n  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建记录时间', \n  `create_by` varchar(32) NOT NULL DEFAULT '' COMMENT '创建人', \n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间', \n  `update_by` varchar(32) NOT NULL DEFAULT '' COMMENT '最后修改人', \n  PRIMARY KEY (`id`), \n  UNIQUE KEY `uk_owner_id` (`owner_id`) \n) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='配置-认领人表'; \n ----------------------------------------\n 略...... \n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SQLFormat(tt.args.sql), "SQLFormat(%v)", tt.args.sql)
		})
	}
}
