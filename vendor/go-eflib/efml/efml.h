/////////////////////////////////////////////////////////////////////////////
//
//  @file efml.h
//
//  @brief API interface of Enflame Management Library
//
//  Enflame Tech, All Rights Reserved. 2024-2025 Copyright (C)
//
/////////////////////////////////////////////////////////////////////////////

#ifndef LIB_EFML_INCLUDE_EFML_H_
#define LIB_EFML_INCLUDE_EFML_H_

#include <stdbool.h>
#include <stdint.h>

#include "efml_export.h"
#ifdef __cplusplus
extern "C" {
#endif

typedef enum {
    EFML_SUCCESS                        = 0,  // Success
    EFML_ERROR_UNINITIALIZED            = 1,  // Error since uninitialized.
    EFML_ERROR_INVALID_ARGUMENT         = 2,
    EFML_ERROR_NOT_SUPPORTED            = 3,
    EFML_ERROR_LIBRARY_NOT_FOUND        = 4,
    EFML_ERROR_INVALID_ERROR_CODE       = 5,
    EFML_ERROR_DRIVER_NOT_LOADED        = 6,
    EFML_ERROR_GCU_LARE_PORT_NUMBER_ERR = 7,
    EFML_ERROR_INVALID_INPUT            = 8,
    EFML_ERROR_FUNCTION_NOT_FOUND       = 9,
    EFML_ERROR_OPEN_DRIVER_VERSION      = 10,
    EFML_ERROR_DRIVER_NOT_COMPATIBLE    = 11,
    EFML_ERROR_NO_DEVICE                = 12,
    EFML_ERROR_DEPRECATION              = 13,
    EFML_ERROR_TIMEOUT                  = 253,
    EFML_ERROR_FAIL                     = 254,
    EFML_ERROR_MAX                      = 255,
} efmlReturn_t;

typedef enum {
    EFML_LINK_SPEED_GEN1 = 1,
    EFML_LINK_SPEED_GEN2 = 2,
    EFML_LINK_SPEED_GEN3 = 3,
    EFML_LINK_SPEED_GEN4 = 4,
    EFML_LINK_SPEED_GEN5 = 5,
} efmlPcieSpeed_t;

typedef enum {
    GCU200 = 0,
    GCU210 = 1,
    GCU300 = 2,
} efmlHwArch_t;

typedef enum {
    UNKNOWN  = 0,
    GCU_ARCH_200  = 20, /* pavo */
    GCU_ARCH_210  = 21, /* dorado */
    GCU_ARCH_300  = 30, /* scorpio */
    GCU_ARCH_400  = 40, /* libra */
} efmlHwArchV2_t;

typedef enum {
    CABLE_NOT_QUALIFIED = 0,
    CABLE_QUALIFIED     = 1,
} efmlCableQualifySts_t;

typedef enum {
    EFML_LINK_WIDTH_X1  = 1,
    EFML_LINK_WIDTH_X2  = 2,
    EFML_LINK_WIDTH_X4  = 4,
    EFML_LINK_WIDTH_X8  = 8,
    EFML_LINK_WIDTH_X16 = 16,
} efmlPcieWidth_t;

typedef enum {
    EFML_GCU_LARE_LINK_SPEED_GEN1      = 1,
    EFML_GCU_LARE_LINK_SPEED_GEN2      = 2,
    EFML_GCU_LARE_LINK_SPEED_GEN3      = 3,
    EFML_GCU_LARE_LINK_SPEED_GEN4      = 4,
    EFML_GCU_LARE_LINK_SPEED_GEN5      = 5,
    EFML_GCU_LARE_LINK_SPEED_ESM_2P5GT = 0x10,
    EFML_GCU_LARE_LINK_SPEED_ESM_5GT,
    EFML_GCU_LARE_LINK_SPEED_ESM_8GT,
    EFML_GCU_LARE_LINK_SPEED_ESM_16GT,
    EFML_GCU_LARE_LINK_SPEED_ESM_20GT,
    EFML_GCU_LARE_LINK_SPEED_ESM_25GT
} efmlGcuLareSpeed_t;

typedef enum {
    EFML_GCU_LARE_LINK_WIDTH_X1  = 1,
    EFML_GCU_LARE_LINK_WIDTH_X2  = 2,
    EFML_GCU_LARE_LINK_WIDTH_X4  = 4,
    EFML_GCU_LARE_LINK_WIDTH_X8  = 8,
    EFML_GCU_LARE_LINK_WIDTH_X16 = 16,
} efmlGcuLareWidth_t;

typedef enum {
    EFML_GCU_LARE_PORT_RC = 0,
    EFML_GCU_LARE_PORT_EP = 1,
} efmlGcuLarePortType_t;

typedef enum {
    EFML_HBM_SCAN_INVALID = 0,
    EFML_HBM_SCAN_START   = 1,
} efmlHbmScanType_t;

typedef enum {
    EFML_EVENT_UNKNOWN          = 0,
    EFML_EVENT_DTU_SUSPEND      = 3,
    EFML_EVENT_DTU_RESET_START  = 10,
    EFML_EVENT_DTU_RESET_FINISH = 11,
    EFML_EVENT_DTU_EID          = 20,
} efmlEventType_t;

typedef enum {
    EFML_GCU_VIRT_MODE_HOST,
    EFML_GCU_VIRT_MODE_VGCU,
    EFML_GCU_VIRT_MODE_MDEV,
    EFML_GCU_VIRT_MODE_SRIOV_PF,
    EFML_GCU_VIRT_MODE_SRIOV_VF,
} efmlGcuVirtMode_t;

typedef enum {
    EFML_MCM_SINGLE_MODE = 1,
    EFML_MCM_DUAL_MODE   = 2,
} efmlMcmMode_t;

typedef enum {
    EFML_ACTIVE  = 1,
    EFML_SUSPEND = 5,
    EFML_SLEEP   = 7,
} efmlPowerMode_t;

#define MAX_CHAR_BUFF_LEN (128)
#define MAX_DEVICE_COUNT (128)
#define MAX_PROCESS_COUNT (64)

/**
 * @brief Enflame data structure for device info.
 *
 */
typedef struct {
    char     name[MAX_CHAR_BUFF_LEN];
    uint32_t vendor_id;
    uint32_t device_id;
    uint32_t domain_id;
    uint32_t bus_id;
    uint32_t dev_id;
    uint32_t func_id;
    uint32_t logic_id;
} efmlDeviceInfo_t;

/**
* @brief Enflame data structure for device power info.
*
*/
typedef struct {
    float pwr_capability;
    float cur_pwr_consumption;
} efmlDevPowerInfo_t;

/**
* @brief Enflame data structure for device memory info.
*
*/
typedef struct {
    uint64_t mem_total_size;
    uint32_t mem_used;
} efmlDevMemInfo_t;

typedef struct {
    uint32_t        event_id;
    efmlEventType_t event_type;
    char            event_msg[MAX_CHAR_BUFF_LEN];
} efmlEvent_t;

/**
* @brief Enflame data structure for cluster memory info.
*
*/
typedef struct {
    uint64_t mem_total_size;
    uint32_t mem_used;
} efmlClusterHbmMemInfo_t;

/**
* @brief Enflame data structure for device clock info.
*
*/
typedef struct {
    uint32_t cur_hbm_clock;
    uint32_t cur_dtu_clock;
} efmlDevClkInfo_t;

typedef struct {
    float cur_asic_temp;
    float cur_mem_temp;
    float cur_board_temp;
} efmlDevThermalInfoV2_t;

/**
* @brief Enflame data structure for device volt info.
*
*/
typedef struct {
    float vdd_dtu;
    float vdd_soc;
    float vdd_hbmqc;
    float vdd_1v8;
    float vdd_vddp;
} efmlDevVoltInfo_t;

/**
 * @brief Enflame data structure for Pcie link info.
 *
 */
typedef struct {
    efmlPcieSpeed_t link_speed;
    efmlPcieSpeed_t max_link_speed;
    efmlPcieWidth_t link_width;
    efmlPcieWidth_t max_link_width;
} efmlPcieLinkInfo_t;

/**
* @brief Enflame data structure for Pcie Throughput info.
*
*/
typedef struct {
    float    tx_throughput;
    float    rx_throughput;
    uint64_t tx_nak;
    uint64_t rx_nak;
} efmlPcieThroughputInfo_t;

/**
* @brief Enflame data structure for GCU-LARE link info.
*
*/
typedef struct {
    efmlGcuLareSpeed_t link_speed;
    efmlGcuLareSpeed_t max_link_speed;
    efmlGcuLareWidth_t link_width;
    efmlGcuLareWidth_t max_link_width;
} efmlGcuLareLinkInfo_t;

/**
*@brief Enflame data structure for GCU-LARE Throughput info.
*
*/
typedef struct {
    float    tx_throughput;
    float    rx_throughput;
    uint64_t tx_nak;
    uint64_t rx_nak;
} efmlGcuLareThroughputInfo_t;

/**
 * @brief Enflame data structure for GCU-LARE port info.
 *
 */
typedef struct {
    uint32_t          connected;
    char              uuid[16];
    uint32_t          vendor_id;
    uint32_t          device_id;
    uint32_t          domain_id;
    uint32_t          bus_id;
    uint32_t          dev_id;
    uint32_t          func_id;
    uint32_t          port_id;
    efmlGcuLarePortType_t port_type;

    uint32_t          remote_card_id;
    char              remote_uuid[16];
    uint32_t          remote_vendor_id;
    uint32_t          remote_device_id;
    uint32_t          remote_domain_id;
    uint32_t          remote_bus_id;
    uint32_t          remote_dev_id;
    uint32_t          remote_func_id;
    uint32_t          remote_port_id;
    efmlGcuLarePortType_t remote_port_type;
} efmlGcuLarePortInfo_t;

typedef struct {
    bool     enabled;  // true: enabled               false: disabled
    bool     pending;  // true: pending               false: no pending
    bool     pdblack;  // true: pending blacklist yes false: no
    uint32_t ecnt_sb;  // single bit error count
    uint32_t ecnt_db;  // double bit error count
} efmlEccStatus_t;

/**
*@brief Enflame data structure for Rma status.
*
*/
typedef struct {
    bool is_dtu_support_rma;
    bool flags;  // true: hint, false: not hint
} efmlRmaStatus_t;

/**
*@brief Enflame data structure for Rma details.
*
*/
typedef struct {
    bool     is_dtu_support_rma;
    bool     flags;  // true: hint, false: not hint
    uint32_t dbe_count;
} efmlRmaDetails_t;

/**
*@brief Enflame data structure for Decoder info.
*
*/
typedef struct {
    uint32_t decoder_inst_num;    // decoder instance number
    float    decoder_resolution;  // decoder resolution
    float    decoder_fps;         // decoder fps
} efmlDecoderCap_t;

/**
*@brief Enflame data structure for HBM usage info.
*
*/
typedef struct {
    uint64_t hbm_total;
    uint64_t hbm_free;
    uint64_t hbm_used;
} efmlHbmUsage_t;

/**
*@brief Enflame data structure for SIP usage info.
*
*/
typedef struct {
    uint32_t sip_total;
    uint32_t sip_init;
    uint32_t sip_idle;
    uint32_t sip_busy;
    uint32_t sip_masked;
    uint32_t sip_hwerr;
} efmlSipUsage_t;

/**
*@brief Enflame data structure for CQM usage info.
*
*/
typedef struct {
    uint32_t cqm_total;
    uint32_t cqm_init;
    uint32_t cqm_idle;
    uint32_t cqm_busy;
    uint32_t cqm_masked;
    uint32_t cqm_hwerr;
} efmlCqmUsage_t;

typedef struct {
    uint32_t pid;           /* process pid         */
    uint64_t dev_mem_usage; /* device memory usage */
    uint64_t sys_mem_usage; /* system memory usage */
} efmlProcessInfo_t;

typedef enum {
      OP_SPACE_GCU,
      OP_SPACE_PF,
      OP_SPACE_CTL,
} efmlOperateSpace_t;

typedef struct {
    uint32_t encoder_usage;
    uint32_t decoder_usage;
} efmlVpuUsage_t;

typedef enum efmlResetType_enum
{
    EFML_RESET_TYPE_NONE = 0,
    EFML_RESET_TYPE_HOT  = 1,
    EFML_RESET_TYPE_FLR  = 2,
    EFML_RESET_TYPE_BLR  = 3,
    EFML_RESET_TYPE_QLR  = 4,
} efmlResetType_t;


typedef enum efmlResetSrc_enum
{
    EFML_RESET_SRC_USER     = 0,
    EFML_RESET_SRC_INTERNAL = 1,
    EFML_RESET_SRC_SIP      = 2,
    EFML_RESET_SRC_BUS      = 3,
    EFML_RESET_SRC_FW       = 4,
    EFML_RESET_SRC_DTE      = 5,
    EFML_RESET_SRC_DRAM_HBM = 6,
    EFML_RESET_SRC_PCIE     = 7,
    EFML_RESET_SRC_UNKNOWN  = 999,
} efmlResetSrc_t;


typedef struct efmlErrorRecord_st
{
    unsigned long long event_id;
    efmlResetType_t    reset_type;
    efmlResetSrc_t     source;
    unsigned long long time_stamp;
    unsigned int       bdf;
} efmlErrorRecord_t;

/**
 * @brief Enflame Management Library initialization
 *
 * @return efmlReturn_t
 */
EFML_EXPORT efmlReturn_t EfmlInit(bool no_driver);

/**
 * @brief Enflame Management Library Shutdown.
 *
 */
EFML_EXPORT void EfmlShutdown();

/**
 * @brief Enflame Management Library map error code to string.
 *
 */
EFML_EXPORT efmlReturn_t EfmlErrorString(efmlReturn_t result, char *p_error_str);

/**
 * @brief Enflame Management Library get driver version info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDriverVer(char *p_driver_ver);

/**
 * @brief Enflame Management Library get efml lib version info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetLibVer(char *p_self_ver);

/**
 * @brief Enflame Management Library get enflame driver status and config access path.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDriverAccessPoint(char *p_enflame_driver_ap);

/**
 * @brief Get the total number of supported devices
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevCount(uint32_t *dev_count);

/**
 * @brief Get the total number of enabled virtual devices of a physical device
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetVdevCount(uint32_t dev_idx, uint32_t *vdev_count);

/**
 * @brief Get the virtual devices index in operate system
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetVdevList(uint32_t dev_idx, uint32_t *vdev_ids, uint32_t *count);

/**
 * @brief Get the total number of maximum supported virtual devices of a physical device
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetMaxVdevCount(uint32_t dev_idx, uint32_t *vdev_count);

/**
 * @brief Enflame Management Library get the device name.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevName(uint32_t dev_idx, char *p_name);

/**
 * @brief Enflame Management Library get the device temperature.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevTempV2(uint32_t dev_idx, efmlDevThermalInfoV2_t *p_temp);

/**
 * @brief Enflame Management Library get the device voltage.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevVolt(uint32_t dev_idx, efmlDevVoltInfo_t *p_volt);

/**
 * @brief Enflame Management Library get the device current power consumption.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevPwr(uint32_t dev_idx, efmlDevPowerInfo_t *p_pwr);

/**
 * @brief Enflame Management Library get the device mem info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevMem(uint32_t dev_idx, efmlDevMemInfo_t *p_mem);

/**
 * @brief Enflame Management Library get the device usage.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevDtuUsage(uint32_t dev_idx, float *p_data);

/**
 * @brief Enflame Management Library get the virtual device mem info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetVdevMem(uint32_t dev_idx, uint32_t vdev_idx, efmlDevMemInfo_t *p_mem);

/**
 * @brief Enflame Management Library get the virtual device usage.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetVdevDtuUsage(uint32_t dev_idx, uint32_t vdev_idx, float *p_data);

EFML_EXPORT efmlReturn_t EfmlGetDevSipUsage(uint32_t dev_idx, float *p_data);

/**
 * @brief Enflame Management Library get the device usage from background sample thread.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevDtuUsageAsync(uint32_t dev_idx, float *p_data);

/**
 * @brief Enflame Management Library get if dtu is in low power mode.
 *
 */
EFML_DEPRECATED_EXPORT efmlReturn_t EfmlGetDevIsLowPowerMode(uint32_t dev_idx, bool *is_low_power_mode);

/**
 * @brief Enflame Management Library get if dtu is health.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevIsHealth(uint32_t dev_idx, bool *is_health);

/**
 * @brief Enflame Management Library get the device max clock frequency.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetMaxFreq(uint32_t dev_idx, uint32_t *max_freq_mhz);

/**
 * @brief Enflame Management Library set the device max clock frequency.
 *
 */
EFML_EXPORT efmlReturn_t EfmlSetMaxFreq(uint32_t dev_idx, uint32_t max_freq_mhz);

/**
 * @brief Enflame Management Library get the device clock info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevClk(uint32_t dev_idx, efmlDevClkInfo_t *p_clk);

/**
 * @brief Enflame Management Library get the memory clock info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetMemClk(uint32_t dev_idx, uint32_t *p_clk_mhz);

/**
 * @brief Enflame Management Library get device info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevInfo(uint32_t dev_idx, efmlDeviceInfo_t *p_info);

/**
 * @brief Enflame Management Library display device topology.
 *
 */
EFML_EXPORT efmlReturn_t EfmlDisplayDevTop(uint32_t dev_idx);

/**
 * @brief Enflame Management Library get firmware version info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetFwVersion(uint32_t dev_idx, char *fw_ver);

/**
 * @brief Enflame Management Library get device UUID info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevUuid(uint32_t dev_idx, char *p_info);

/**
 * @brief Enflame Management Library get device SN(Serial Number) info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevSn(uint32_t dev_idx, char *p_sn);

/**
 * @brief Enflame Management Library get device PN(Part Number) info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevPn(uint32_t dev_idx, char *p_pn);

/**
 * @brief Enflame Management Library get device Manufacturing Date.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevMfd(uint32_t dev_idx, char *p_date);

/**
 * @brief Enflame Management Library get device SKU info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevSKU(uint32_t dev_idx, char *p_info);

/**
 * @brief Enflame Management Library get device hw arch info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetHwArch(uint32_t dev_idx, efmlHwArch_t *p_arch);

/**
 * @brief Enflame Management Library get device hw arch name.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetHwArchName(uint32_t dev_idx, char *arch_name);

/**
 * @brief Enflame Management Library get power cable qualify status.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetCableQualify(
    uint32_t dev_idx, efmlCableQualifySts_t *p_sts);

/**
 * @brief Enflame Management Library get device PCIe slot or OAM name.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevSlotOamName(uint32_t dev_idx, char *p_slot_oam);

/**
 * @brief Enflame Management Library get current device pcie link speed.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetPcieLinkSpeed(uint32_t dev_idx, efmlPcieSpeed_t *p_link_speed);

/**
 * @brief Enflame Management Library get current device pcie link width.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetPcieLinkWidth(uint32_t dev_idx, efmlPcieWidth_t *p_link_width);

/**
 * @brief Enflame Management Library get current device pcie link info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetPcieLinkInfo(uint32_t dev_idx, efmlPcieLinkInfo_t *p_link_info);

/**
 * @brief Enflame Management Library get pcie throughput info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetPcieThroughput(uint32_t dev_idx, efmlPcieThroughputInfo_t *p_info);

/**
 * @brief Enflame Management Library pcie hot reset.
 *
 */
EFML_EXPORT efmlReturn_t EfmlPcieHotReset(uint32_t dev_idx);

/**
 * @brief Enflame Management Library advance pcie hot reset.
 *
 */
EFML_EXPORT efmlReturn_t EfmlPcieHotResetV2(uint32_t dev_idx, bool is_force);

/**
 * @brief Enflame Management Library system pcie reset.
 *
 */
EFML_EXPORT efmlReturn_t EfmlPcieHotResetV3(uint32_t dev_idx);

/**
 * @brief Enflame Management Library advance pcie function level reset.
 *
 */
EFML_EXPORT efmlReturn_t EfmlPcieFLR(uint32_t dev_idx, bool is_force);

/**
 * @brief Enflame Management Library get pcie slot ID.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetPciePhysicalSlotID(uint32_t dev_idx, uint32_t *p_id);

/**
 * @brief Enflame Management Library get dram ECC status.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetEccStatus(uint32_t dev_idx, uint32_t *status);

/**
 * @brief Enflame Management Library get total GCU-LARE port numbers.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetGcuLarePortNum(uint32_t dev_idx, uint32_t *p_data);

/**
 * @brief Enflame Management Library get GCU-LARE port info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetGcuLarePortInfo(uint32_t dev_idx, uint32_t port_id, efmlGcuLarePortInfo_t *p_info);

/**
 * @brief Enflame Management Library get GCU-LARE link info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetGcuLareLinkInfo(uint32_t dev_idx, uint32_t port_id, efmlGcuLareLinkInfo_t *p_info);

/**
 * @brief Enflame Management Library get GCU-LARE dtuid info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetGcuLareId(uint32_t dev_idx, uint32_t *p_data);

/**
 * @brief Enflame Management Library get GCU-LARE support info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetGcuLareIsSupported(uint32_t dev_idx, bool *is_gcu_lare_supported);

/**
 * @brief Enflame Management Library get GCU-LARE throughput info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetGcuLareThroughput(
    uint32_t dev_idx, uint32_t port_id, efmlGcuLareThroughputInfo_t *p_info);

/**
 * @brief Enflame Management Library get the total number of special device's pgs.
 *
 */
EFML_DEPRECATED_EXPORT efmlReturn_t EfmlGetPGCount(uint32_t dev_idx, uint32_t *pg_count);

/**
 * @brief Enflame Management Library get the device pg usage.
 *
 */
EFML_DEPRECATED_EXPORT efmlReturn_t EfmlGetDevPGUsage(uint32_t dev_idx, uint32_t pg_idx, float *p_data);

/**
 * @brief Enflame Management Library get the device pg usage from background sample thread..
 *
 */
EFML_DEPRECATED_EXPORT efmlReturn_t EfmlGetDevPGUsageAsync(uint32_t dev_idx, uint32_t pg_idx, float *p_data);

/**
 * @brief Enflame Management Library get the total number of special device's clusters.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetClusterCount(uint32_t dev_idx, uint32_t *cluster_count);

/**
 * @brief Enflame Management Library get the device cluster usage.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevClusterUsage(uint32_t dev_idx, uint32_t cluster_idx, float *p_data);

/**
 * @brief Enflame Management Library get the cluster mem info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevClusterHbmMem(
    uint32_t dev_idx, uint32_t cluster_idx, efmlClusterHbmMemInfo_t *p_mem);

/**
 * @brief Enflame Management Library get the device health stauts.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevHealth(uint32_t dev_idx, bool *health);

/**
 * @brief Enflame Management Library get the device ecc stauts.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevEccStatus(uint32_t dev_idx, efmlEccStatus_t *p_status);

/**
 * @brief Enflame Management Library get the device rma stauts.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevRmaStatus(uint32_t dev_idx, efmlRmaStatus_t *p_status);

/**
 * @brief Enflame Management Library get the device rma details.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevRmaDetails(uint32_t dev_idx, efmlRmaDetails_t *p_details);

/**
 * @brief Enflame Management Library get the device logic id.
 *
 */
EFML_DEPRECATED_EXPORT efmlReturn_t EfmlGetDevLogicId(uint32_t dev_idx, uint32_t *p_logic_id);

/**
 * @brief Enflame Management Library get the device decoder capability.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetDevDecoderCap(uint32_t dev_idx, efmlDecoderCap_t *p_decoder_cap);

/**
 * @brief Enflame Management Library set the device ecc mode.
 *
 */
EFML_EXPORT efmlReturn_t EfmlSetDevEccMode(uint32_t dev_idx, bool enable);

/**
 * @brief Enflame Management Library set the device hbm scan mode.
 *
 */
EFML_EXPORT efmlReturn_t EfmlHbmScanMode(uint32_t dev_idx, efmlHbmScanType_t op_type);

/**
 * @brief Enflame Management Library get event message.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetEvent(int timeout_ms, efmlEvent_t *p_event);

/**
 * @brief Enflame Management Library get event message. Non-blocking.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetEventV2(uint32_t dev_idx, efmlEvent_t *p_event);

/**
 * @brief Enflame Management Library start listen event.
 *
 */
EFML_DEPRECATED_EXPORT efmlReturn_t EfmlStartListenEvent(uint32_t dev_idx);

/**
 * @brief Enflame Management Library switch operate space.
 *
 */
EFML_EXPORT efmlReturn_t EfmlSwitchOperateSpace(efmlOperateSpace_t opsp);

/**
 * @brief Enflame Management Library get process info on device.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetProcessInfo(
    uint32_t dev_idx, uint32_t *process_count, efmlProcessInfo_t *p_info);

/**
 * @brief Enflame Management Library get device numa node.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetNumaNode(uint32_t dev_idx, int *numa_node);

/**
 * @brief Enflame Management Library get PF name of device. Only support on VF.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetSriovPfName(uint32_t dev_idx, char *pf_name);

/**
 * @brief Enflame Management Library get device affinity CPU list.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetAffinityCpuList(uint32_t dev_idx, char *cpu_list);

/**
 * @brief Enflame Management Library get SR-IOV status.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetGcuVirtStatus(uint32_t dev_idx, efmlGcuVirtMode_t *virt_mode);

/**
 * @brief Enflame Management Library get SR-IOV VF number. Only support on PF.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetSriovVfNum(uint32_t dev_idx, uint32_t *num);


/**
 * @brief Enflame Management Library get VPU usage.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetVpuUsage(uint32_t dev_idx, efmlVpuUsage_t *p_usage);

/**
 * @brief Enflame Management Library get multi chip modules mode.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetMcmMode(uint32_t dev_idx, efmlMcmMode_t *mcm_mode);

/**
 * @brief Enflame Management Library set multi chip modules mode. Only support on PF.
 *
 */
EFML_EXPORT efmlReturn_t EfmlSetMcmMode(uint32_t dev_idx, efmlMcmMode_t mcm_mode);


/**
 * @brief Enflame Management Library get power mode.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetPowerMode(uint32_t dev_idx, efmlPowerMode_t *power_mode);


/**
 * @brief Enflame Management Library set power mode.
 *
 */
EFML_EXPORT efmlReturn_t EfmlSetPowerMode(uint32_t dev_idx, efmlPowerMode_t power_mode);


/**
 * @brief Enflame Management Library get RAS info.
 *
 */
EFML_EXPORT efmlReturn_t EfmlGetErrorRecords(uint32_t dev_idx, uint32_t count, efmlErrorRecord_t *p_records,
    uint32_t *total_count);


/////////////////////////////////////////////////////////////////////////////
//
// efml next
//
/////////////////////////////////////////////////////////////////////////////

#define EFML_STRUCT_VERSION(data, ver) (unsigned int)(sizeof(efml ## data ## _v ## ver ## _t) | \
                                                      (ver << 24U))

typedef struct efmlDevice* efmlDevice_t;

/**
 * Return values for EFML API calls.
 */
typedef enum efmlReturn_enum
{
    // cppcheck-suppress *
    EFML_SUCCESS_V2 = 0,                       //!< The operation was successful
    EFML_ERROR_UNINITIALIZED_V2 = 1,           //!< EFML was not first initialized with efmlInit()
    EFML_ERROR_INVALID_ARGUMENT_V2 = 2,        //!< A supplied argument is invalid
    EFML_ERROR_NOT_SUPPORTED_V2 = 3,           //!< The requested operation is not available on target device
    EFML_ERROR_NO_PERMISSION = 4,              //!< The current user does not have permission for operation
    EFML_ERROR_NOT_FOUND = 6,                  //!< A query to find an object was unsuccessful
    EFML_ERROR_INSUFFICIENT_SIZE = 7,          //!< An input argument is not large enough
    EFML_ERROR_INSUFFICIENT_POWER = 8,         //!< A device's external power cables are not properly attached
    EFML_ERROR_DRIVER_NOT_LOADED_V2 = 9,       //!< ENFLAME driver is not loaded
    EFML_ERROR_TIMEOUT_V2 = 10,                //!< User provided timeout passed
    EFML_ERROR_IRQ_ISSUE = 11,                 //!< ENFLAME Kernel detected an interrupt issue with a GCU
    EFML_ERROR_LIBRARY_NOT_FOUND_V2 = 12,      //!< EFML Shared Library couldn't be found or loaded
    EFML_ERROR_FUNCTION_NOT_FOUND_V2 = 13,     //!< Local version of EFML doesn't implement this function
    EFML_ERROR_CORRUPTED_INFOROM = 14,         //!< infoROM is corrupted
    EFML_ERROR_GCU_IS_LOST = 15,               //!< The GCU has fallen off the bus or has otherwise become inaccessible
    EFML_ERROR_RESET_REQUIRED = 16,            //!< The GCU requires a reset before it can be used again
    EFML_ERROR_OPERATING_SYSTEM = 17,          //!< The GCU control device has been blocked by the operating system/cgroups
    EFML_ERROR_LIB_RM_VERSION_MISMATCH = 18,   //!< RM detects a driver/library version mismatch
    EFML_ERROR_IN_USE = 19,                    //!< An operation cannot be performed because the GCU is currently in use
    EFML_ERROR_MEMORY = 20,                    //!< Insufficient memory
    EFML_ERROR_NO_DATA = 21,                   //!< No data
    EFML_ERROR_VGCU_ECC_NOT_SUPPORTED = 22,    //!< The requested vgpu operation is not available on target device, because ECC is enabled
    EFML_ERROR_INSUFFICIENT_RESOURCES = 23,    //!< Ran out of critical resources, other than memory
    EFML_ERROR_FREQ_NOT_SUPPORTED = 24,        //!< Ran out of critical resources, other than memory
    EFML_ERROR_ARGUMENT_VERSION_MISMATCH = 25, //!< The provided version is invalid/unsupported
    EFML_ERROR_DEPRECATED  = 26,               //!< The requested functionality has been deprecated
    EFML_ERROR_NOT_READY = 27,                 //!< The system is not ready for the request
    EFML_ERROR_GCU_NOT_FOUND = 28,             //!< No GCUs were found
    EFML_ERROR_INVALID_STATE = 29,             //!< Resource not in correct state to perform requested operation
    EFML_ERROR_UNKNOWN = 999                   //!< An internal driver error occurred
} efmlReturn_v1_t;


/**
 * Simplified chip architecture
 */
#define EFML_DEVICE_ARCH_GCU300 30 // Devices based on the Enflame Scorpio architecture
#define EFML_DEVICE_ARCH_GCU400 40 // Devices based on the Enflame Libra architecture

#define EFML_DEVICE_ARCH_UNKNOWN   0xffffffff // Anything else, presumably something newer

typedef unsigned int efmlDeviceArchitecture_t;


/**
 * Generic enable/disable enum.
 */
typedef enum efmlEnableState_enum
{
    EFML_FEATURE_DISABLED    = 0,     //!< Feature disabled
    EFML_FEATURE_ENABLED     = 1      //!< Feature enabled
} efmlEnableState_t;


/**
 * Clock types.
 *
 * All speeds are in Mhz.
 */
typedef enum efmlClockType_enum
{
    EFML_CLOCK_GRAPHICS  = 0,        //!< Graphics clock domain
    EFML_CLOCK_SIP       = 1,        //!< SIP clock domain
    EFML_CLOCK_MEM       = 2,        //!< Memory clock domain
    EFML_CLOCK_VIDEO     = 3,        //!< Video encoder/decoder clock domain

    // Keep this last
    EFML_CLOCK_COUNT //!< Count of clock types
} efmlClockType_t;

/**
 * Memory allocation information for a device.
 */
typedef struct efmlMemory_st
{
    unsigned long long total;        //!< Total installed device memory (in bytes)
    unsigned long long free;         //!< Unallocated device memory (in bytes)
    unsigned long long used;         //!< Allocated device memory (in bytes).
} efmlMemory_t;


/**
 * Memory allocation information for a device.
 */
typedef struct efmlMemory_v2_st
{
    unsigned int version;            //!< Structure format version (must be 2)
    unsigned long long total;        //!< Total installed device memory (in bytes)
    unsigned long long free;         //!< Unallocated device memory (in bytes)
    unsigned long long used;         //!< Allocated device memory (in bytes).
    unsigned long long reserved;     //!< Device memory (in bytes) reserved for system use (driver or firmware).
} efmlMemory_v2_t;


/**
 * Information about running compute processes on the GPU
 */
typedef struct efmlProcessInfo_st
{
    unsigned int pid;                 //!< Process ID
    unsigned long long usedGcuMemory; //!< Amount of used GCU memory in bytes.
    unsigned int gcuInstanceId;       //!< always 0xFFFFFFFF
    unsigned int computeInstanceId;   //!< always 0xFFFFFFFF
} efmlProcessInfo_v1_t;


/**
 * Structure to store utilization value and process Id
 */
typedef struct efmlProcessUtilizationSample_st
{
    unsigned int pid;                   //!< PID of process
    unsigned long long timeStamp;       //!< CPU Timestamp in microseconds
    unsigned int sipUtil;               //!< SIP Util Value
    unsigned int memUtil;               //!< Memory Util Value
    unsigned int encUtil;               //!< Encoder Util Value. Not support now
    unsigned int decUtil;               //!< Decoder Util Value. Not support now
} efmlProcessUtilizationSample_t;


/**
 * Buffer size guaranteed to be large enough for pci bus id
 */
#define EFML_DEVICE_PCI_BUS_ID_BUFFER_SIZE      32

/**
 * PCI information about a GPU device.
 */
typedef struct efmlPciInfo_st
{
    unsigned int domain;             //!< The PCI domain on which the device's bus resides, 0 to 0xffffffff
    unsigned int bus;                //!< The bus on which the device resides, 0 to 0xff
    unsigned int device;             //!< The device's id on the bus, 0 to 31
    unsigned int function;           //!< The device's function number, 0 to 7
    unsigned int pciDeviceId;        //!< The combined 16-bit device id and 16-bit vendor id

    unsigned int pciSubSystemId;     //!< The 32-bit Sub System Device ID

    char busId[EFML_DEVICE_PCI_BUS_ID_BUFFER_SIZE]; //!< The tuple domain:bus:device.function PCI identifier (&amp; NULL terminator)
} efmlPciInfo_t;

/**
 * Represents the queryable PCIe utilization counters
 */
typedef enum efmlPcieUtilCounter_enum
{
    EFML_PCIE_UTIL_TX_BYTES             = 0, // 1KB granularity
    EFML_PCIE_UTIL_RX_BYTES             = 1, // 1KB granularity

    // Keep this last
    EFML_PCIE_UTIL_COUNT
} efmlPcieUtilCounter_t;

/**
 * Temperature thresholds.
 */
typedef enum efmlTemperatureThresholds_enum
{
    EFML_TEMPERATURE_THRESHOLD_SHUTDOWN = 0,    // Temperature at which the GCU will shut down
                                                // for HW protection
    EFML_TEMPERATURE_THRESHOLD_SLOWDOWN = 1,    // Temperature at which the GCU will slowdown to the minimum
                                                // clock
    EFML_TEMPERATURE_THRESHOLD_MEM_MAX  = 2,    // Max temperature at which the memory works properly
    EFML_TEMPERATURE_THRESHOLD_GCU_MAX  = 3,    // Max temperature at which the GCU core works properly
    // Keep this last
    EFML_TEMPERATURE_THRESHOLD_COUNT
} efmlTemperatureThresholds_t;

/**
 * Temperature sensors.
 */
typedef enum efmlTemperatureSensors_enum
{
    EFML_TEMPERATURE_GCU      = 0,    //!< Temperature sensor for the GCU die

    // Keep this last
    EFML_TEMPERATURE_COUNT
} efmlTemperatureSensors_t;


/*!
 * GCU virtualization mode types.
 */
typedef enum efmlGcuVirtualizationMode_enum {
    EFML_GCU_VIRTUALIZATION_MODE_NONE = 0,         //!< Represents Bare Metal GCU
    EFML_GCU_VIRTUALIZATION_MODE_PASSTHROUGH = 1,  //!< Device is associated with GCU-Passthorugh
    EFML_GCU_VIRTUALIZATION_MODE_VGCU = 2,         //!< Device is associated with vGCU.
    EFML_GCU_VIRTUALIZATION_MODE_MDEV = 3,         //!< Device is associated with MDEV.
    EFML_GCU_VIRTUALIZATION_MODE_SRIOV_PF = 4,     //!< Device is associated with SRIOV PF.
    EFML_GCU_VIRTUALIZATION_MODE_SRIOV_VF = 5,     //!< Device is associated with SRIOV VF.
} efmlGcuVirtualizationMode_t;


typedef struct efmlEGcuLareInfo_st {
    unsigned int maxTthroughput;  //!< The max throughput in Gbps
    unsigned long long mac[8];    //!< The MAC address of the GCU-LARE
    unsigned int ip[8];           //!< The IP address of the GCU-LARE
    unsigned char remoteUuid[16]; //!< The UUID of the remote GCU-LARE
    char nicName[16];             //!< The NIC name of the GCU-LARE
    unsigned long long throughput;
} efmlGcuLareInfo_t;


/**
 * Utilization information for a device.
 */
typedef struct efmlUtilization_st
{
    unsigned int gcu;                //!< Percent of time over the past sample period during which one or more kernels was executing on the GCU
    unsigned int memory;             //!< Percent of time over the past sample period during which global (device) memory was being read or written
} efmlUtilization_t;


typedef struct efmlOccupancy_st
{
    unsigned int sip;
    unsigned int sip_vector;
    unsigned int sip_matrix;
    unsigned int sip_special_function;
    unsigned int memory;
} efmlOccupancy_t;


/**
 * @brief Initialize the EFML library.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 *
 * @note This function must be called before any other EFML functions.
 */
EFML_EXPORT efmlReturn_v1_t efmlInit(void);
EFML_EXPORT efmlReturn_v1_t efmlInit_v2(void);


/**
 * @brief Initialize the EFML library for physical functions.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 *
 * @note This function must be called before any other EFML functions.
 * @note Only call this function if SR-IOV is enabled.
 * @note This function is conflict with efmlInit(). If efmlInit() is called, efmlPfInit() will do nothing.
 */
EFML_EXPORT efmlReturn_v1_t efmlPfInit(void);


/**
 * @brief Shutdown the EFML library.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 *
 * @note This function must be called to release resources allocated by efmlInit().
 *      After calling this function, no other EFML functions should be called.
 */
EFML_EXPORT efmlReturn_v1_t efmlShutdown(void);


/**
 * @brief Shutdown the EFML library for physical functions.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 *
 * @note This function must be called to release resources allocated by efmlInit().
 *      After calling this function, no other EFML functions should be called.
 * @note Only call this function after efmlPfInit().
 */
EFML_EXPORT efmlReturn_v1_t efmlPfShutdown(void);


/**
 * @brief Helper method for converting EFML error codes to strings.
 *
 * @param error [in] The error code to convert.
 *
 * @return A string representation of the error code.
 */
EFML_EXPORT const char* efmlErrorString(efmlReturn_v1_t error);

/**
 * @brief Get architecture for device
 *
 * @param device [in] The device handle.
 * @param arch [out] Reference where architecture is returned.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetArchitecture(efmlDevice_t device, efmlDeviceArchitecture_t* arch);

/**
 * @brief Retrieves the current clock speeds for the device.
 *
 * @param device [in] The device handle.
 * @param type [in] The type of clock to retrieve.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetClockInfo(efmlDevice_t device, efmlClockType_t type, unsigned int* clock);

efmlReturn_v1_t efmlDeviceGetComputeRunningProcesses_v3(
    efmlDevice_t device, unsigned int* infoCount, efmlProcessInfo_v1_t* info);


/**
 * @brief Get the total number of devices.
 *
 * @param count [out] The total number of devices.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetCount_v2(unsigned int* count);


/**
 * @brief Retrieves the current PCIe link generation
 *
 * @param device [in] The device handle.
 * @param currLinkGen [out] The current PCIe link generation.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetCurrPcieLinkGeneration(efmlDevice_t device, unsigned int* currLinkGen);


/**
 * @brief Retrieves the current PCIe link width
 *
 * @param device [in] The device handle.
 * @param currLinkWidth [out] The current PCIe link width.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetCurrPcieLinkWidth(efmlDevice_t device, unsigned int* currLinkWidth);

/**
 * @brief Retrieves the current occupancy for the decoder.
 *
 * @param device [in] The device handle.
 * @param occupancy [out] The current occupancy for the decoder.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetDecoderOccupancy(efmlDevice_t device, unsigned int* occupancy);

/**
 * @brief Retrieves the current utilization and sampling size in microseconds for the decoder.
 *
 * @param device [in] The device handle.
 * @param utilization [out] The current utilization for the decoder.
 * @param samplingPeriodUs [out] The sampling period in microseconds for the utilization counter.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetDecoderUtilization(efmlDevice_t device, unsigned int* utilization,
    unsigned int* samplingPeriodUs);

/**
 * @brief Retrieves the current occupancy for the encoder.
 *
 * @param device [in] The device handle.
 * @param occupancy [out] The current occupancy for the encoder.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetEncoderOccupancy(efmlDevice_t device, unsigned int* occupancy);

/**
 * @brief Retrieves the current utilization and sampling size in microseconds for the encoder.
 *
 * @param device [in] The device handle.
 * @param utilization [out] The current utilization for the encoder.
 * @param samplingPeriodUs [out] The sampling period in microseconds for the utilization counter.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetEncoderUtilization(efmlDevice_t device, unsigned int* utilization,
    unsigned int* samplingPeriodUs);

/**
 * @brief Get the effective power limit enforced by the driver.
 *
 * @param device [in] The device handle.
 * @param limit [out] The effective power limit in milliwatts.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetEnforcedPowerLimit(efmlDevice_t device, unsigned int* limit);

EFML_EXPORT efmlReturn_v1_t efmlDeviceGetFanSpeed(efmlDevice_t device, unsigned int* speed);
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetFanSpeed_v2(efmlDevice_t device, unsigned int fan, unsigned int* speed);

EFML_EXPORT efmlReturn_v1_t efmlDeviceGetGraphicsRunningProcesses_v3(
    efmlDevice_t device, unsigned int* infoCount, efmlProcessInfo_v1_t* info);

/**
 * @brief Get the device handle by index.
 *
 * @param index [in] The index of the device to get.
 * @param device [out] The device handle to return.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetHandleByIndex_v2(unsigned int index, efmlDevice_t* device);

EFML_EXPORT efmlReturn_v1_t efmlDeviceGetMPSComputeRunningProcesses_v3(
    efmlDevice_t device, unsigned int* infoCount, efmlProcessInfo_v1_t* info);


/**
 * @brief Acquire the handle for a particular device, based on its PCI bus id.
 *
 * @param pciBusId [in] The PCI bus id of the device to get.
 * @param device [out] The device handle to return.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetHandleByPciBusId_v2(const char* pciBusId, efmlDevice_t* device);


/**
 * @brief Retrieves the maximum clock speeds for the device.
 *
 * @param device [in] The device handle.
 * @param type [in] The type of clock to retrieve.
 * @param clock [out] The maximum clock speed for the device in MHz.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetMaxClockInfo(efmlDevice_t device, efmlClockType_t type, unsigned int* clock);


/**
 * @brief Retrieves the maximum PCIe link generation possible with this device and system
 *
 * @param device [in] The device handle.
 * @param maxLinkGen [out] The maximum PCIe link generation.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetMaxPcieLinkGeneration(efmlDevice_t device, unsigned int* maxLinkGen);


/**
 * @brief Retrieves the maximum PCIe link width possible with this device and system
 *
 * @param device [in] The device handle.
 * @param maxLinkWidth [out] The maximum PCIe link width.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetMaxPcieLinkWidth(efmlDevice_t device, unsigned int* maxLinkWidth);


/**
 * @brief Retrieves the amount of used, free, reserved and total memory available on the device, in bytes.
 *        Enabling ECC reduces the amount of total available memory, due to the extra required parity bits.
 *
 * @param device [in] The device handle.
 * @param memory [out] The memory information for the device.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetMemoryInfo(efmlDevice_t device, efmlMemory_t* memory);
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetMemoryInfo_v2(efmlDevice_t device, efmlMemory_v2_t* memory);


/**
 * @brief Retrieves the device's name
 *
 * @param device [in] The device handle.
 * @param name [out] The device's name.
 * @param length [in] The length of the name buffer.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetName(efmlDevice_t device, char* name, unsigned int length);


/**
 * @brief Retrieves the device's PCI information
 *
 * @param device [in] The device handle.
 * @param pciInfo [out] The device's PCI information.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetPciInfo_v3(efmlDevice_t device, efmlPciInfo_t* pciInfo);


/**
 * @brief Retrieve PCIe utilization information. This function is querying a byte counter over a 1ms interval and thus
 *        is the PCIe throughput over that interval.
 *
 * @param device [in] The device handle.
 * @param counter [in] The PCIe utilization counter to query.
 * @param value [out] Reference in which to return throughput in KB/s
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetPcieThroughput(
    efmlDevice_t device, efmlPcieUtilCounter_t counter, unsigned int* value);

/**
 * @brief Retrieves power usage for this GCU in milliwatts
 *
 * @param device [in] The device handle.
 * @param power [out] The power usage in milliwatts
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetPowerUsage(efmlDevice_t device, unsigned int* power);

/**
 * @brief Retrieves the current utilization and process ID
 *
 * @param device [in] The device handle.
 * @param utilization [out] Pointer to caller-supplied buffer in which guest process utilization samples are returned
 * @param processSamplesCount [in/out] Pointer to caller-supplied array size, and returns number of processes running
 * @param lastSeenTimeStamp [in] Return only samples with timestamp greater than lastSeenTimeStamp.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetProcessUtilization(efmlDevice_t device,
    efmlProcessUtilizationSample_t* utilization, unsigned int* processSamplesCount,
    unsigned long long lastSeenTimeStamp);

efmlReturn_v1_t efmlDeviceTestGcuLareConnection(efmlDevice_t device,
    const char * source_address, const char * dest_address, int timeout_in_s, bool * result);


/**
 * @brief Retrieves the state of the device's GCU-LARE for the port specified
 *
 * @param device [in] The device handle.
 * @param port [in] The port id to query.
 * @param isActive [out] The state of the GCU-LARE for the port specified.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
efmlReturn_v1_t efmlDeviceGetGcuLareState(efmlDevice_t device, unsigned int port, efmlEnableState_t *isActive);


/**
 * @brief Retrieves the GCU-LARE information for the port specified
 *
 * @param device [in] The device handle.
 * @param port [in] The port id to query.
 * @param info [out] The GCU-LARE information for the port specified.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
efmlReturn_v1_t efmlDeviceGetGcuLareInfo(efmlDevice_t device, unsigned int port, efmlGcuLareInfo_t *info);


efmlReturn_v1_t efmlDeviceGetGcuLareRemotePciInfo_v2(efmlDevice_t device, unsigned int link, efmlPciInfo_t *pciInfo);

/**
 * @brief Retrieves the current temperature readings for the device, in degrees C.
 *
 * @param device [in] The device handle.
 * @param sensorType [in] The type of temperature sensor to retrieve.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetTemperature(efmlDevice_t device, efmlTemperatureSensors_t sensorType,
    int* temperature);

/**
 * Retrieves the temperature threshold for the GCU with the specified threshold type in degrees C.
 *
 * @param device [in] The device handle.
 * @param thresholdType [in] The type of temperature threshold to retrieve.
 * @param temperature [out] The temperature threshold in degrees C.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetTemperatureThreshold(efmlDevice_t device,
    efmlTemperatureThresholds_t thresholdType, unsigned int* temperature);


/**
 * @brief Retrieves the current occupancy for the device's major subsystems.
 *
 * @param device [in] The device handle.
 * @param occupancy [out] The device's occupancy rates.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetOccupancy(efmlDevice_t device, efmlOccupancy_t* occupancy);

/**
 * @brief Retrieves the current utilization rates for the device's major subsystems.
 *        The utilization rates are returned as a percentage of time over the past sample period during which
 *        one or more kernels was executing on the GCU, and global (device) memory was being read or written.
 *
 * @param device [in] The device handle.
 * @param utilization [out] The device's utilization rates.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetUtilizationRates(efmlDevice_t device, efmlUtilization_t* utilization);

/**
 * @brief Retrieves the device's virtualization mode
 *
 * @param device [in] The device handle.
 * @param mode [out] The device's virtualization mode.
 *
 * @return EFML_SUCCESS_V2 on success, or an error code on failure.
 */
EFML_EXPORT efmlReturn_v1_t efmlDeviceGetVirtualizationMode(efmlDevice_t device, efmlGcuVirtualizationMode_t* mode);

EFML_EXPORT efmlReturn_v1_t efmlDeviceGetMigMode(efmlDevice_t device, unsigned int* currentMode,
    unsigned int* pendingMode);

#ifdef __cplusplus
}  // extern "C"
#endif

#include "efml_legacy.h"

#endif  // LIB_EFML_INCLUDE_EFML_H_

